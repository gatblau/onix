/*
    Onix CMDB - Copyright (c) 2018-2019 by www.gatblau.org

    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
    Unless required by applicable law or agreed to in writing, software distributed under
    the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
    either express or implied.
    See the License for the specific language governing permissions and limitations under the License.

    Contributors to this project, hereby assign copyright in this code to the project,
    to be licensed under the same terms as the rest of the code.

*/
DO
$$
BEGIN

/*
  Encapsulates the logic to determine the status of a record update:
    - N: no update as no changes found - new and old records are the same
    - L: no update as the old record was updated by another client before this update could be committed
    - U: update - the record was updated successfully
 */
CREATE OR REPLACE FUNCTION get_update_status(
    current_version bigint, -- the version of the record in the database
    local_version bigint, -- the version in the new specified record
    updated boolean -- whether or not the record was updated in the database by the last update statement
  )
  RETURNS char(1)
  LANGUAGE 'plpgsql'
  COST 100
  VOLATILE
AS $BODY$
  DECLARE
    result char(1);
  BEGIN
    -- if there were not rows affected
    IF NOT updated THEN
      -- if the local version is the same as the record version
      IF (local_version = current_version) THEN
        -- no update was required as required record was the same as stored record
        result := 'N';
      ELSE
        -- no update was made as stored record is optimistically locked
        -- i.e. updated by other client before this update can be committed
        result := 'L';
      END IF;
    ELSE
      -- the stored record was updated successfully
      result := 'U';
    END IF;

    RETURN result;
  END;
$BODY$;

ALTER FUNCTION get_update_status(bigint, bigint, boolean)
OWNER TO onix;

/*
  Validates that an item attribute store contains the keys required or allowed
  by its item type definition.
 */
CREATE OR REPLACE FUNCTION validate_item_attr(item_key_param character varying)
  RETURNS VOID
  LANGUAGE 'plpgsql'
  COST 100
  VOLATILE
AS $BODY$
DECLARE
  item_type_key character varying (100);
  rule record;
  store hstore;
BEGIN
  -- gets the item type key for the specified item
  SELECT item_type.key
  FROM item
    INNER JOIN item_type
      ON item.item_type_id = item_type.id
  WHERE item.key = item_key_param
  INTO item_type_key;

  -- gets the attribute hstore in the item
  SELECT attribute FROM item WHERE key = item_key_param INTO store;

  -- loop through the item type attr_valid hstore key-value pairs
  FOR rule IN
    SELECT (each(attr_valid)).*
    FROM item_type
    WHERE key = item_type_key
  LOOP
    IF rule.value = 'required' THEN
      IF NOT (store ? rule.key) THEN
        RAISE EXCEPTION 'Key % is required', rule.key
          USING HINT = 'Check the key exist in the item attribute store.';
      END IF;
    END IF;
  END LOOP;
END;
$BODY$;

---------------------------------------------------------
-- UPSERT ITEM
-- a) inserts a new item if it is not there and return true; or,
-- b) updates an existing item, if any of their values is different from the value in the database and returns false; or,
-- c) does not perform any update if a record exists and the passed in values are the same as the values in the database
---------------------------------------------------------
CREATE OR REPLACE FUNCTION upsert_item(
    key_param character varying,
    name_param character varying,
    description_param text,
    meta_param jsonb,
    tag_param text[],
    attribute_param hstore,
    status_param smallint,
    item_type_key_param character varying,
    local_version_param bigint,
    changedby_param character varying
  )
  RETURNS TABLE(result char(1))
  LANGUAGE 'plpgsql'
  COST 100
  VOLATILE
AS $BODY$
  DECLARE
    result char(1); -- the result status for the upsert
    current_version bigint; -- the version of the row before the update or null if no row
    rows_affected integer;
    item_type_id_value integer;
  BEGIN
    -- find the item type surrogate key from the provided natural key
    SELECT id FROM item_type WHERE key = item_type_key_param INTO item_type_id_value;
    IF (item_type_id_value IS NULL) THEN
      -- the provided natural key is not in the item type table, cannot proceed
     RAISE EXCEPTION 'Nonexistent Item Type Key --> %', item_type_key_param
        USING HINT = 'Check an Item Type with the key exist in the database.';
    END IF;

    -- get the item current version
    SELECT version FROM item WHERE key = key_param INTO current_version;
    -- if no version is found then go for an insert
    IF (current_version IS NULL) THEN
      INSERT INTO item (
        id,
        key,
        name,
        description,
        meta,
        tag,
        attribute,
        status,
        item_type_id,
        version,
        created,
        updated,
        changedby
      )
      VALUES (
          nextval('item_id_seq'),
          key_param,
          name_param,
          description_param,
          meta_param,
          tag_param,
          attribute_param,
          status_param,
          item_type_id_value,
          1,
          current_timestamp,
          null,
          changedby_param
      );
      result := 'I';
    ELSE
      -- if a version is found, go for an update
      UPDATE item SET
        name = name_param,
        description = description_param,
        meta = meta_param,
        tag = tag_param,
        attribute = attribute_param,
        status = status_param,
        item_type_id = item_type_id_value,
        version = version + 1,
        updated = current_timestamp,
        changedby = changedby_param
      WHERE key = key_param
      -- the database record has not been modified by someone else
      AND local_version_param = current_version
      AND (
        -- the fields to be updated have not changed
        name != name_param OR
        description != description_param OR
        status != status_param OR
        item_type_id != item_type_id_value OR
        meta != meta_param OR
        tag != tag_param OR
        attribute != attribute_param
      );
      -- determines if the update has gone ahead
      GET DIAGNOSTICS rows_affected := ROW_COUNT;
      -- works out the update status
      SELECT get_update_status(current_version, local_version_param, rows_affected > 0) INTO result;
    END IF;
    RETURN QUERY SELECT result;
  END;
  $BODY$;

ALTER FUNCTION upsert_item(character varying,character varying,text,jsonb, text[],hstore,smallint,character varying, bigint, character varying)
OWNER TO onix;

---------------------------------------------------------
-- UPSERT ITEM_TYPE
-- a) inserts a new item type if it is not there and return true; or,
-- b) updates an existing item type, if any of their values is different from the value in the database and returns false; or,
-- c) does not perform any update if a record exists and the passed in values are the same as the values in the database
---------------------------------------------------------
CREATE OR REPLACE FUNCTION upsert_item_type(
    key_param character varying,
    name_param character varying,
    description_param text,
    attr_valid_param hstore, -- keys allowed or required in item attributes
    local_version_param bigint,
    changedby_param character varying
  )
  RETURNS TABLE(result char(1))
  LANGUAGE 'plpgsql'
  COST 100
  VOLATILE
AS $BODY$
  DECLARE
    result char(1); -- the result status for the upsert
    current_version bigint; -- the version of the row before the update or null if no row
    rows_affected integer;
BEGIN
  SELECT version FROM item_type WHERE key = key_param INTO current_version;
  IF (current_version IS NULL) THEN
    INSERT INTO item_type (
      id,
      key,
      name,
      description,
      attr_valid,
      custom,
      version,
      created,
      updated,
      changedby
    )
    VALUES (
      nextval('item_type_id_seq'),
      key_param,
      name_param,
      description_param,
      attr_valid_param,
      true,
      1,
      current_timestamp,
      null,
      changedby_param
    );
    result := 'I';
  ELSE
    UPDATE item_type SET
      name = name_param,
      description = description_param,
      attr_valid = attr_valid_param,
      version = version + 1,
      updated = current_timestamp,
      changedby = changedby_param
    WHERE key = key_param
    AND local_version_param = current_version -- optimistic locking
    AND (
      name != name_param OR
      description != description_param OR
      allowed_keys != allowed_keys_param
    );
    GET DIAGNOSTICS rows_affected := ROW_COUNT;
    SELECT get_update_status(current_version, local_version_param, rows_affected > 0) INTO result;
  END IF;
  RETURN QUERY SELECT result;
END;
$BODY$;

ALTER FUNCTION upsert_item_type(character varying, character varying, text, hstore, bigint, character varying)
OWNER TO onix;

---------------------------------------------------------
-- UPSERT LINK
-- a) inserts a new link between items if it is not there and return true; or,
-- b) updates an existing link, if any of their values is different from the value in the database and returns false; or,
-- c) does not perform any update if a record exists and the passed in values are the same as the values in the database
---------------------------------------------------------
CREATE OR REPLACE FUNCTION upsert_link(
    key_param character varying,
    start_item_key_param character varying,
    end_item_key_param character varying,
    description_param text,
    meta_param jsonb,
    tag_param text[],
    attribute_param hstore,
    local_version_param bigint,
    changedby_param character varying
  )
  RETURNS TABLE(result char(1))
  LANGUAGE 'plpgsql'
  COST 100
  VOLATILE
AS $BODY$
  DECLARE
    result char(1); -- the result status for the upsert
    current_version bigint; -- the version of the row before the update or null if no row
    rows_affected integer;
    start_item_id bigint;
    end_item_id bigint;
BEGIN
  SELECT id FROM item WHERE key = start_item_key_param INTO start_item_id;
  SELECT id FROM item WHERE key = end_item_key_param INTO end_item_id;
  SELECT version FROM link WHERE key = key_param INTO current_version;
  IF (current_version IS NULL) THEN
    INSERT INTO link (
      id,
      key,
      start_item_id,
      end_item_id,
      description,
      meta,
      tag,
      attribute,
      version,
      created,
      updated,
      changedby
    )
    VALUES (
      nextval('link_id_seq'),
      key_param,
      start_item_id,
      end_item_id,
      description_param,
      meta_param,
      tag_param,
      attribute_param,
      1,
      current_timestamp,
      null,
      changedby_param
    );
    result := 'I';
  ELSE
    UPDATE link SET
      meta = meta_param,
      description = description_param,
      tag = tag_param,
      attribute = attribute_param,
      version = version + 1,
      updated = current_timestamp,
      changedby = changedby_param
    WHERE key = key_param
    AND local_version_param = current_version -- optimistic locking
    AND (
      meta != meta_param OR
      description != description_param OR
      tag != tag_param OR
      attribute != attribute_param
    );
    GET DIAGNOSTICS rows_affected := ROW_COUNT;
    SELECT get_update_status(current_version, local_version_param, rows_affected > 0) INTO result;
  END IF;
  RETURN QUERY SELECT result;
END;
$BODY$;

ALTER FUNCTION upsert_link(character varying, character varying, character varying, text, jsonb, text[], hstore, bigint, character varying)
OWNER TO onix;

/*
  gets an item by its natural key.
  use: select * from item('the_item_key')
 */
CREATE OR REPLACE FUNCTION item(key_param character varying)
  RETURNS TABLE(
    id bigint,
    key character varying,
    name character varying,
    description text,
    status smallint,
    item_type_id integer,
    meta jsonb,
    tag text[],
    attribute hstore,
    version bigint,
    created timestamp(6) with time zone,
    updated timestamp(6) with time zone,
    changedby character varying
  )
  LANGUAGE 'plpgsql'
  COST 100
  VOLATILE
AS $BODY$
  BEGIN
    RETURN QUERY SELECT
      i.id,
      i.key,
      i.name,
      i.description,
      i.status,
      i.item_type_id,
      i.meta,
      i.tag,
      i.attribute,
      i.version,
      i.created,
      i.updated,
      i.changedby
    FROM item i
    WHERE i.key = key_param;
  END;
$BODY$;

ALTER FUNCTION item(character varying)
OWNER TO onix;

END;
$$