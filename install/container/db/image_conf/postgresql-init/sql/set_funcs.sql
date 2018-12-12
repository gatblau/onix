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
DO $$
BEGIN

---------------------------------------------------------
-- SET ITEM
---------------------------------------------------------
CREATE OR REPLACE FUNCTION set_item(
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
     RAISE EXCEPTION 'Item Type Key --> % not found.', item_type_key_param
        USING hint = 'Check an Item Type with the key exist in the database.';
    END IF;

    -- checks that the attributes passed in comply with the validation in the item_type
    PERFORM check_item_attr(item_type_key_param, attribute_param);

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

ALTER FUNCTION set_item(character varying,character varying,text,jsonb, text[],hstore,smallint,character varying, bigint, character varying)
OWNER TO onix;

---------------------------------------------------------
-- SET ITEM_TYPE
---------------------------------------------------------
CREATE OR REPLACE FUNCTION set_item_type(
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
  -- checks that the attribute store parameter contain the correct values
  PERFORM check_attr_valid(attr_valid_param);

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

ALTER FUNCTION set_item_type(character varying, character varying, text, hstore, bigint, character varying)
OWNER TO onix;

---------------------------------------------------------
-- SET LINK_TYPE
---------------------------------------------------------
CREATE OR REPLACE FUNCTION set_link_type(
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
  -- checks that the attribute store parameter contain the correct values
  PERFORM check_attr_valid(attr_valid_param);

  -- gets the link type current version
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

ALTER FUNCTION set_link_type(character varying, character varying, text, hstore, bigint, character varying)
OWNER TO onix;

---------------------------------------------------------
-- SET LINK
---------------------------------------------------------
CREATE OR REPLACE FUNCTION set_link(
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
    RAISE EXCEPTION 'Link Type Key --> % not found.', link_type_key_param
      USING hint = 'Check a Link Type with the key exist in the database.';
  END IF;

  -- checks that the attributes passed in comply with the validation in the link_type
  PERFORM check_link_attr(link_type_key_param, attribute_param);

  SELECT id FROM item WHERE key = start_item_key_param INTO start_item_id;
  IF (start_item_id IS NULL) THEN
    -- the start item does not exist
    RAISE EXCEPTION 'Start item with key --> % does not exist.', start_item_key_param
      USING hint = 'Check an item with the specified key exist in the database.';
  END IF;

  SELECT id FROM item WHERE key = end_item_key_param INTO end_item_id;
  IF (end_item_id IS NULL) THEN
    -- the end item does not exist
    RAISE EXCEPTION 'End item with key --> % does not exist.', end_item_key_param
      USING hint = 'Check an item with the specified key exist in the database.';
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

ALTER FUNCTION set_link(character varying, character varying, character varying, character varying, text, jsonb, text[], hstore, bigint, character varying)
  OWNER TO onix;

END
$$;