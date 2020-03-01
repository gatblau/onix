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
     ox_get_enc_key_usage: gets the number of items and links that are using a specified encryption key
        i.e. key 1 or key 2 respectively.
        Use this query to understand the state of key rotation at a point in time.
     */
    CREATE OR REPLACE FUNCTION ox_get_enc_key_usage(
        enc_key_ix_param smallint,
        role_key_param character varying[]
    )
    RETURNS BIGINT
    LANGUAGE 'plpgsql'
    COST 100
    STABLE
    AS $BODY$
        DECLARE
            item_count bigint;
            link_count bigint;
    BEGIN
        SELECT count(*) INTO item_count
         FROM item i
          INNER JOIN partition p on i.partition_id = p.id
          INNER JOIN privilege pr on p.id = pr.partition_id
          INNER JOIN role r on pr.role_id = r.id
         WHERE i.enc_key_ix = enc_key_ix_param
           AND pr.can_read = TRUE
           AND r.key = ANY(role_key_param);

        SELECT count(*) INTO link_count
        FROM link l
            INNER JOIN item i on l.start_item_id = i.id -- gets the partition from the start item is linked to
             INNER JOIN partition p on i.partition_id = p.id
             INNER JOIN privilege pr on p.id = pr.partition_id
             INNER JOIN role r on pr.role_id = r.id
        WHERE i.enc_key_ix = enc_key_ix_param
          AND pr.can_read = TRUE
          AND r.key = ANY(role_key_param);

        RETURN (item_count + link_count)::BIGINT;
    END
    $BODY$;

    ALTER FUNCTION ox_get_enc_key_usage(smallint, character varying[])
        OWNER TO onix;

END
$$;