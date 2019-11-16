/*
    Onix Config Manager - Copyright (c) 2018-2019 by www.gatblau.org

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
    /*
      checks that the specified role can modify partition, privilege and role tables
      based on its role level
     */
    CREATE OR REPLACE FUNCTION ox_can_manage_partition(
      role_key_param character varying[]
    )
      RETURNS VOID
      LANGUAGE 'plpgsql'
      COST 100
      STABLE
    AS
    $BODY$
    DECLARE
      level integer;
      owner character varying(100);
    BEGIN
      -- finds the logged role level and owner
      SELECT r.level, r.owner
      FROM role r
      WHERE r.key = ANY(role_key_param)
      ORDER BY r.level DESC
      LIMIT 1
        INTO level, owner;

      IF (level = 0) THEN
        RAISE EXCEPTION 'Role % is not authorised to modify role/partition information.', role_key_param
          USING hint = 'The role is a level 0 role, it needs to be level 1 or 2.';
      ELSEIF (level = 1 AND owner != ANY(role_key_param)) THEN
        RAISE EXCEPTION 'Role % is not authorised to modify a role/partition with a different owner.', role_key_param
          USING hint = 'The role is a level 1 role, it needs to be level 2.';
      END IF;
    END;
    $BODY$;

    ALTER FUNCTION ox_can_manage_partition(character varying[])
      OWNER TO onix;

    /*
      Encapsulates the logic to determine the status of a record update:
        - N: no update as no changes found - new and old records are the same
        - L: no update as the old record was updated by another client before this update could be committed
        - U: update - the record was updated successfully
     */
    CREATE OR REPLACE FUNCTION ox_get_update_status(
      current_version bigint, -- the version of the record in the database
      local_version bigint, -- the version in the new specified record
      updated boolean -- whether or not the record was updated in the database by the last update statement
    )
      RETURNS char(1)
      LANGUAGE 'plpgsql'
      COST 100
      STABLE
    AS
    $BODY$
    DECLARE
      result char(1);
    BEGIN
      -- if there were not rows affected
      IF NOT updated THEN
        -- if the local version is the same as the record version
        IF (local_version = current_version OR local_version IS NULL) THEN
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

    ALTER FUNCTION ox_get_update_status(bigint, bigint, boolean)
      OWNER TO onix;

    /*
      Validates that an item attribute store contains the keys required or allowed
      by its item type definition.
     */
    CREATE OR REPLACE FUNCTION ox_check_item_attr(item_type_key character varying, attributes hstore)
      RETURNS VOID
      LANGUAGE 'plpgsql'
      COST 100
      STABLE
    AS
    $BODY$
    DECLARE
      validation_rules hstore;
      rule             record;
    BEGIN
      -- gets the validation rules for an item attributes
      SELECT attr_valid INTO validation_rules
      FROM item_type
      WHERE key = item_type_key;

      IF NOT FOUND THEN
        RAISE EXCEPTION 'Invalid item type key ''%''.', item_type_key
          USING hint = 'Check an item type with such key has been defined.';
      END IF;

      -- if validation is defined at the item_type then
      -- validate the item attribute field
      IF NOT (validation_rules IS NULL) THEN
        -- loop through the validation rules key-regex pairs
        -- to validate 'required' key compliance in the passed-in attributes
        FOR rule IN
          SELECT (each(validation_rules)).*
          LOOP
            IF (rule.value = 'required') THEN
              IF NOT (attributes ? rule.key) THEN
                RAISE EXCEPTION 'Item of type ''%'' requires attribute ''%''.', item_type_key, rule.key
                  USING hint =
                      'Where required attributes are specified in the item type, a request to insert or update an item of that type must also specify the regex of the required attribute(s).';
              END IF;
            END IF;
          END LOOP;

        -- loop through the passed-in item attribute hstore key-regex pairs
        -- to validate 'allowed' key compliance in the passed-in attributes
        FOR rule IN
          SELECT (each(attributes)).*
          LOOP
            IF NOT (validation_rules ? rule.key) THEN
              RAISE EXCEPTION 'Attribute ''%'' is not allowed!', rule.key
                USING hint = 'Revise the item attributes removing the attribute not allowed.';
            END IF;
          END LOOP;
      END IF;
    END;
    $BODY$;

    ALTER FUNCTION ox_check_item_attr(character varying, hstore)
      OWNER TO onix;

    /*
      Validates that a link attribute store contains the keys required or allowed
      by its link type definition.
     */
    CREATE OR REPLACE FUNCTION ox_ox_check_link_attr(link_type_key character varying, attributes hstore)
      RETURNS VOID
      LANGUAGE 'plpgsql'
      COST 100
      STABLE
    AS
    $BODY$
    DECLARE
      validation_rules hstore;
      rule             record;
    BEGIN
      -- gets the validation rules for a link attributes
      SELECT attr_valid INTO validation_rules
      FROM link_type
      WHERE key = link_type_key;

      IF NOT FOUND THEN
        RAISE EXCEPTION 'Invalid link type key ''%''.', link_type_key
          USING hint = 'Check a link type with such a key has been defined.';
      END IF;

      -- if validation is defined at the item_type then
      -- validate the item attribute field
      IF NOT (validation_rules IS NULL) THEN
        -- loop through the validation rules key-regex pairs
        -- to validate 'required' key compliance in the passed-in attributes
        FOR rule IN
          SELECT (each(validation_rules)).*
          LOOP
            IF (rule.value = 'required') THEN
              IF NOT (attributes ? rule.key) THEN
                RAISE EXCEPTION 'Attribute ''%'' is required and was not provided.', rule.key
                  USING hint =
                      'Where required attributes are specified in the link type, a request to insert or update a link of that type must also specify the regex of the required attribute(s).';
              END IF;
            END IF;
          END LOOP;

        -- loop through the passed-in item attribute hstore key-regex pairs
        -- to validate 'allowed' key compliance in the passed-in attributes
        FOR rule IN
          SELECT (each(attributes)).*
          LOOP
            IF NOT (validation_rules ? rule.key) THEN
              RAISE EXCEPTION 'Attribute ''%'' is not allowed!', rule.key
                USING hint = 'Revise the item attributes removing the attribute not allowed.';
            END IF;
          END LOOP;
      END IF;
    END;
    $BODY$;

    ALTER FUNCTION ox_ox_check_link_attr(character varying, hstore)
      OWNER TO onix;

    /*
      checks that the specified hstore contains only 'required' or 'allowed' values
     */
    CREATE OR REPLACE FUNCTION ox_check_attr_valid(attributes hstore)
      RETURNS VOID
      LANGUAGE 'plpgsql'
      COST 100
      STABLE
    AS
    $BODY$
    DECLARE
      rule record;
    BEGIN
      -- if the attributes hstore is defined
      IF NOT (attributes IS NULL) THEN
        -- loop through the validation rules key-regex pairs
        -- to determine if there are values other than 'required' or 'allowed'
        FOR rule IN
          SELECT (each(attributes)).*
          LOOP
            IF NOT ((rule.value = 'required') OR (rule.value = 'allowed')) THEN
              RAISE EXCEPTION 'Attribute ''%'' has an invalid regex: ''%''.', rule.key, rule.value
                USING hint = 'Attribute values can only be either ''required'' or ''allowed''';
            END IF;
          END LOOP;
      END IF;
    END;
    $BODY$;

    ALTER FUNCTION ox_check_attr_valid(hstore)
      OWNER TO onix;

    /*
      select ox_check_link(link_type_key_param, start_item_type_key_param, end_item_type_key_param):
        checks that a link of a given type is valid (i.e. can be used to join to items of given types
        in a particular direction)
     */
    CREATE OR REPLACE FUNCTION ox_check_link(link_type_key_param character varying,
                                          start_item_type_key_param character varying,
                                          end_item_type_key_param character varying)
      RETURNS VOID
      LANGUAGE 'plpgsql'
      COST 100
      STABLE
    AS
    $BODY$
    DECLARE
      rule_count integer;
    BEGIN
      SELECT COUNT(*) INTO rule_count
      FROM link_rule r
             INNER JOIN link_type lt
                        ON lt.id = r.link_type_id
             INNER JOIN item_type start_item_type
                        ON r.start_item_type_id = start_item_type.id
             INNER JOIN item_type end_item_type
                        ON r.end_item_type_id = end_item_type.id
      WHERE lt.key = link_type_key_param
        AND start_item_type.key = start_item_type_key_param
        AND end_item_type.key = end_item_type_key_param;

      IF (rule_count = 0) THEN
        RAISE EXCEPTION 'Unallowed link: a link of type ''%'' cannot be used to connect from items of type ''%'' to items of type ''%''.', link_type_key_param, start_item_type_key_param, end_item_type_key_param
          USING hint = 'Check the link type is correct and the direction of the link is allowed.';
      END IF;
    END;
    $BODY$;

    ALTER FUNCTION ox_check_link(character varying, character varying, character varying)
      OWNER TO onix;

  END
  $$;