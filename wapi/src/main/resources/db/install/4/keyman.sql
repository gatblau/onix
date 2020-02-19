/*
    Onix Config Manager - Copyright (c) 2018-2020 by www.gatblau.org

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
     key management queries
     */

    /*
     ox_get_enc_key_usage: gets the number of items that have encrypted meta and/or txt attributes and whether they
        are using key 1 or key 2 respectively.
        Use this query to understand the state of key rotation at a point in time.
     */
    CREATE OR REPLACE FUNCTION ox_get_enc_key_usage(
        keyno_param int,
        role_key_param character varying[]
    )
    RETURNS TABLE(key_count VARCHAR)
    LANGUAGE 'plpgsql'
    COST 100
    STABLE
    AS $BODY$
        DECLARE key_char VARCHAR;
    BEGIN
        IF keyno_param = 1 THEN
            key_char = '\001';
        ELSIF keyno_param = 2 THEN
            key_char = '\002';
        ELSE
            RAISE 'Invalid key no';
        END IF;
        RETURN QUERY SELECT count(*)::VARCHAR as key_count
         FROM item i
          INNER JOIN partition p on i.partition_id = p.id
          INNER JOIN privilege pr on p.id = pr.partition_id
          INNER JOIN role r on pr.role_id = r.id
         WHERE (substring(i.meta_enc from 4 for 1)::VARCHAR = key_char OR substring(i.txt_enc from 4 for 1)::VARCHAR = key_char)
           AND pr.can_read = TRUE
           AND r.key = ANY(role_key_param);
    END
    $BODY$;

    ALTER FUNCTION ox_get_enc_key_usage(int, character varying[])
        OWNER TO onix;

    /*
     ox_get_enc_items: gets a list of items that have encrypted meta and/or txt fields using a specific key (i.e. 1 or 2)
      Use this query to retrieve items whose encryption keys have to be rotated.
     */
    CREATE OR REPLACE FUNCTION ox_get_enc_items(
        key_no_param int, -- whether to query key 1 or key 2
        max_items_param int, -- cap the result set
        role_key_param character varying[] -- the user role
    )
        RETURNS TABLE
        (
            id bigint,
            key character varying,
            name character varying,
            description text,
            status smallint,
            item_type_key character varying,
            meta jsonb,
            meta_enc bytea,
            txt text,
            txt_enc bytea,
            tag text[],
            attribute hstore,
            version bigint,
            created timestamp(6) with time zone,
            updated timestamp(6) with time zone,
            changed_by character varying,
            model_key character varying,
            partition_key character varying
        )
        LANGUAGE 'plpgsql'
        COST 100
        STABLE
    AS $BODY$
    DECLARE key_code VARCHAR;
    BEGIN
        IF key_no_param = 1 THEN
            key_code = '\001';
        ELSIF key_no_param = 2 THEN
            key_code = '\002';
        ELSE
            RAISE 'key_no_param is out of range. Use 1 or 2.';

        END IF;

        RETURN QUERY SELECT
             i.id,
             i.key,
             i.name,
             i.description,
             i.status,
             it.key as item_type_key,
             i.meta,
             i.meta_enc,
             i.txt,
             i.txt_enc,
             i.tag,
             i.attribute,
             i.version,
             i.created,
             i.updated,
             i.changed_by,
             m.key as model_key,
             p.key as partition_key
         FROM item i
              INNER JOIN item_type it ON i.item_type_id = it.id
              INNER JOIN model m ON m.id = it.model_id
              INNER JOIN partition p on i.partition_id = p.id
              INNER JOIN privilege pr on p.id = pr.partition_id
              INNER JOIN role r on pr.role_id = r.id
         WHERE (substring(i.meta_enc FROM 4 FOR 1)::VARCHAR = key_code OR substring(i.txt_enc FROM 4 FOR 1)::VARCHAR = key_code)
            AND pr.can_read = TRUE
            AND r.key = ANY(role_key_param)
         LIMIT max_items_param;
    END
    $BODY$;

    ALTER FUNCTION ox_get_enc_items(int, int, character varying[])
        OWNER TO onix;
END
$$;