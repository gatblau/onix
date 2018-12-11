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
  STABLE
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
CREATE OR REPLACE FUNCTION validate_item_attr(item_type_key character varying, attributes hstore)
  RETURNS VOID
  LANGUAGE 'plpgsql'
  COST 100
  STABLE
AS $BODY$
DECLARE
  validation_rules hstore;
  rule record;
BEGIN
  -- gets the validation rules for an item attributes
  SELECT attr_valid INTO validation_rules
  FROM item_type
  WHERE key = item_type_key;

  IF NOT FOUND THEN
    RAISE EXCEPTION 'Invalid item type key ''%''.', item_type_key
      USING HINT = 'Check an item type with such key has been defined.';
  END IF;

  -- if validation is defined at the item_type then
  -- validate the item attribute field
  IF NOT (validation_rules IS NULL) THEN
    -- loop through the validation rules key-value pairs
    -- to validate 'required' key compliance in the passed-in attributes
    FOR rule IN
      SELECT (each(validation_rules)).*
    LOOP
      IF (rule.value = 'required') THEN
        IF NOT (attributes ? rule.key) THEN
          RAISE EXCEPTION 'Attribute ''%'' is required and was not provided.', rule.key
            USING HINT = 'Where required attributeS are specified in the item type, a request to insert or update an item of that type must also specify the value of the required attribute(s).';
        END IF;
      END IF;
    END LOOP;

    -- loop through the passed-in item attribute hstore key-value pairs
    -- to validate 'allowed' key compliance in the passed-in attributes
    FOR rule IN
      SELECT (each(attributes)).*
    LOOP
      IF NOT (validation_rules ? rule.key) THEN
        RAISE EXCEPTION 'Attribute ''%'' is not allowed!', rule.key
          USING HINT = 'Revise the item attributes removing the attribute not allowed.';
      END IF;
    END LOOP;
  END IF;
END;
$BODY$;

/*
  Validates that a link attribute store contains the keys required or allowed
  by its link type definition.
 */
CREATE OR REPLACE FUNCTION validate_link_attr(link_type_key character varying, attributes hstore)
  RETURNS VOID
  LANGUAGE 'plpgsql'
  COST 100
  STABLE
AS $BODY$
DECLARE
  validation_rules hstore;
  rule record;
BEGIN
  -- gets the validation rules for a link attributes
  SELECT attr_valid INTO validation_rules
  FROM link_type
  WHERE key = link_type_key;

  IF NOT FOUND THEN
    RAISE EXCEPTION 'Invalid link type key ''%''.', link_type_key
    USING HINT = 'Check a link type with such a key has been defined.';
  END IF;

  -- if validation is defined at the item_type then
  -- validate the item attribute field
  IF NOT (validation_rules IS NULL) THEN
    -- loop through the validation rules key-value pairs
    -- to validate 'required' key compliance in the passed-in attributes
    FOR rule IN
    SELECT (each(validation_rules)).*
    LOOP
      IF (rule.value = 'required') THEN
        IF NOT (attributes ? rule.key) THEN
          RAISE EXCEPTION 'Attribute ''%'' is required and was not provided.', rule.key
          USING HINT = 'Where required attributes are specified in the link type, a request to insert or update a link of that type must also specify the value of the required attribute(s).';
        END IF;
      END IF;
    END LOOP;

    -- loop through the passed-in item attribute hstore key-value pairs
    -- to validate 'allowed' key compliance in the passed-in attributes
    FOR rule IN
    SELECT (each(attributes)).*
    LOOP
        IF NOT (validation_rules ? rule.key) THEN
           RAISE EXCEPTION 'Attribute ''%'' is not allowed!', rule.key
           USING HINT = 'Revise the item attributes removing the attribute not allowed.';
        END IF;
    END LOOP;
  END IF;
END;
$BODY$;

---------------------------------------------------------
-- UPSERT ITEM
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

    -- checks that the attributes passed in comply with the validation in the item_type
    PERFORM validate_item_attr(item_type_key_param, attribute_param);

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
      attr_valid != attr_valid_param
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
-- UPSERT LINK_TYPE
---------------------------------------------------------
CREATE OR REPLACE FUNCTION upsert_link_type(
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
  SELECT version FROM link_type WHERE key = key_param INTO current_version;
  IF (current_version IS NULL) THEN
    INSERT INTO link_type (
      id,
      key,
      name,
      description,
      attr_valid,
      version,
      created,
      updated,
      changedby
    )
    VALUES (
      nextval('link_type_id_seq'),
      key_param,
      name_param,
      description_param,
      attr_valid_param,
      1,
      current_timestamp,
      null,
      changedby_param
    );
    result := 'I';
  ELSE
    UPDATE link_type SET
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
      attr_valid != attr_valid_param
    );
    GET DIAGNOSTICS rows_affected := ROW_COUNT;
    SELECT get_update_status(current_version, local_version_param, rows_affected > 0) INTO result;
  END IF;
  RETURN QUERY SELECT result;
END;
$BODY$;

ALTER FUNCTION upsert_link_type(character varying, character varying, text, hstore, bigint, character varying)
OWNER TO onix;

---------------------------------------------------------
-- UPSERT LINK
---------------------------------------------------------
CREATE OR REPLACE FUNCTION upsert_link(
    key_param character varying,
    link_type_key_param character varying,
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
    link_type_id_value integer;
BEGIN
  -- find the link type surrogate key from the provided natural key
  SELECT id FROM link_type WHERE key = link_type_key_param INTO link_type_id_value;
  IF (link_type_id_value IS NULL) THEN
    -- the provided natural key is not in the link type table, cannot proceed
    RAISE EXCEPTION 'Nonexistent Link Type Key --> %', link_type_key_param
    USING HINT = 'Check a Link Type with the key exist in the database.';
  END IF;

  -- checks that the attributes passed in comply with the validation in the link_type
  PERFORM validate_link_attr(link_type_key_param, attribute_param);

  SELECT id FROM item WHERE key = start_item_key_param INTO start_item_id;
  IF (start_item_id IS NULL) THEN
    -- the start item does not exist
    RAISE EXCEPTION 'Nonexistent start item with key --> %', start_item_key_param
    USING HINT = 'Check an item with the specified key exist in the database.';
  END IF;

  SELECT id FROM item WHERE key = end_item_key_param INTO end_item_id;
  IF (start_item_id IS NULL) THEN
    -- the start item does not exist
    RAISE EXCEPTION 'Nonexistent end item with key --> %', end_item_key_param
    USING HINT = 'Check an item with the specified key exist in the database.';
  END IF;

  SELECT version FROM link WHERE key = key_param INTO current_version;
  IF (current_version IS NULL) THEN
    INSERT INTO link (
      id,
      key,
      link_type_id,
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
      link_type_id_value,
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

ALTER FUNCTION upsert_link(character varying, character varying, character varying, character varying, text, jsonb, text[], hstore, bigint, character varying)
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