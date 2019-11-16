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
      ox_create_tag(...):
        creates a new tag for a subtree starting at the specified root item.
      */
    CREATE OR REPLACE FUNCTION ox_create_tag(root_item_key_param character varying,
                                          tag_label_param character varying,
                                          tag_name_param character varying,
                                          tag_description_param text,
                                          tag_created_by_param character varying)
      RETURNS TABLE
              (
                result char(1)
              )
      LANGUAGE 'plpgsql'
      VOLATILE
    AS
    $BODY$
    DECLARE
      tree_ids     BIGINT[];
      loop_item    RECORD;
      item_store   HSTORE = ''::HSTORE;
      link_store   HSTORE = ''::HSTORE;
      root_item_id BIGINT;
      label_exists BOOLEAN;
      data_exists  BOOLEAN;
    BEGIN
      root_item_id := (SELECT i.id FROM item i WHERE i.key = root_item_key_param);

      tree_ids := (SELECT ox_get_child_items(root_item_id));
      FOR loop_item IN
        SELECT DISTINCT ON (ic.id) ic.id, ic.version
        FROM item_change ic
        WHERE id = ANY (tree_ids::BIGINT[])
        ORDER BY ic.id ASC, ic.changed DESC
        LOOP
          item_store := item_store || hstore(loop_item.id::TEXT, loop_item.version::TEXT);
        END LOOP;

      tree_ids := (SELECT ox_get_child_links(root_item_id));
      FOR loop_item IN
        SELECT DISTINCT ON (lc.id) lc.id, lc.version
        FROM link_change lc
        WHERE id = ANY (tree_ids::BIGINT[])
        ORDER BY lc.id ASC, lc.changed DESC
        LOOP
          link_store := link_store || hstore(loop_item.id::TEXT, loop_item.version::TEXT);
        END LOOP;

      -- checks if label already exists
      SELECT COUNT(*) > 0
      FROM tag
      WHERE label = tag_label_param INTO label_exists;

      -- checks if tag data already exists
      SELECT COUNT(*) > 0
      FROM tag
      WHERE root_item_key = root_item_key_param
        AND item_data = item_store
        AND link_data = link_store INTO data_exists;

      -- only inserts the tag if there is not one with the same label
      IF NOT (label_exists OR data_exists) THEN
        INSERT INTO tag (
           label,
           root_item_key,
           name,
           description,
           item_data,
           link_data,
           version,
           changed_by)
        VALUES (tag_label_param,
                root_item_key_param,
                tag_name_param,
                tag_description_param,
                item_store,
                link_store,
                1,
                tag_created_by_param);
        RETURN QUERY SELECT 'I'::char(1);
      END IF;
      -- if the tag exists then retrieve a conflict
      RETURN QUERY SELECT 'L'::char(1);
    END;
    $BODY$;

    /*
      ox_update_tag(...):
        updates some of the attributes of an existing tag (e.g. label, name and description).
      */
    CREATE OR REPLACE FUNCTION ox_update_tag(root_item_key_param character varying,
                                          current_label_param character varying,
                                          new_label_param character varying,
                                          name_param character varying,
                                          description_param text,
                                          changed_by_param character varying,
                                          local_version_param bigint)
      RETURNS TABLE
              (
                result char(1)
              )
      LANGUAGE 'plpgsql'
      VOLATILE
    AS
    $BODY$
    DECLARE
      result          char(1); -- the result status for the upsert
      current_version bigint; -- the version of the row before the update or null if no row
      rows_affected   integer;
    BEGIN
      -- gets the current version
      SELECT version
      FROM tag
      WHERE root_item_key = root_item_key_param
        AND label = current_label_param INTO current_version;

      UPDATE tag
      SET label       = new_label_param,
          name        = name_param,
          description = description_param,
          version     = version + 1,
          updated     = current_timestamp,
          changed_by  = changed_by_param
      WHERE root_item_key = root_item_key_param
        AND label = current_label_param
        -- concurrency management - optimistic locking (disabled if local_version_param is null)
        AND (local_version_param = current_version OR local_version_param IS NULL)
        AND (
          label != new_label_param OR
          name != name_param OR
          description != description_param
        );
      GET DIAGNOSTICS rows_affected := ROW_COUNT;
      SELECT ox_get_update_status(current_version, local_version_param, rows_affected > 0) INTO result;
      RETURN QUERY SELECT result;
    END;
    $BODY$;

    /*
      ox_delete_tag: deletes the specified tag.
     */
    CREATE OR REPLACE FUNCTION ox_delete_tag(root_item_key_param character varying,
                                          label_param character varying -- if no label, then deletes all tags for the item
    )
      RETURNS VOID
      LANGUAGE 'plpgsql'
      COST 100
      VOLATILE
    AS
    $BODY$
    BEGIN
      DELETE
      FROM tag s
      WHERE (s.label = label_param OR label_param IS NULL)
        AND s.root_item_key = root_item_key_param;
    END
    $BODY$;

    ALTER FUNCTION ox_delete_tag(character varying, character varying)
      OWNER TO onix;

    /*
      ox_get_tree_content(root_item_key_param, label_param): inspects the tag hstores for information
        about a specific tag items and links and retrieve a set of ids and versions for them.
     */
    CREATE OR REPLACE FUNCTION ox_get_tree_content(root_item_key_param character varying,
                                                label_param character varying)
      RETURNS TABLE
              (
                id      text,
                version text,
                is_item boolean
              )
      LANGUAGE 'plpgsql'
      VOLATILE
    AS
    $BODY$
    DECLARE
      item_data HSTORE;
      link_data HSTORE;
    BEGIN
      item_data :=
          (SELECT s.item_data FROM tag s WHERE s.label = label_param AND s.root_item_key = root_item_key_param);
      link_data :=
          (SELECT s.link_data FROM tag s WHERE s.label = label_param AND s.root_item_key = root_item_key_param);
      RETURN QUERY SELECT *, true AS is_item FROM each(item_data);
      RETURN QUERY SELECT *, false AS is_item FROM each(link_data);
    END;
    $BODY$;

    ALTER FUNCTION ox_get_tree_content(character varying, character varying)
      OWNER TO onix;

    /*
      ox_get_tree_items(root_item_key_param, label_param): gets a list of all the items that are part
        of a tag tree for a specific parent item and a label.
    */
    CREATE OR REPLACE FUNCTION ox_get_tree_items(root_item_key_param character varying,
                                              label_param character varying)
      RETURNS TABLE
              (
                operation     character,
                changed       timestamp with time zone,
                id            bigint,
                key           character varying,
                name          character varying,
                description   text,
                meta          jsonb,
                tag           text[],
                attribute     hstore,
                status        smallint,
                item_type_id  integer,
                version       bigint,
                partition_id  bigint,
                created       timestamp with time zone,
                updated       timestamp with time zone,
                changed_by    character varying,
                item_type_key character varying
              )
      LANGUAGE 'plpgsql'
      COST 100
      VOLATILE
      ROWS 1000
    AS
    $BODY$
    BEGIN
      RETURN QUERY
        SELECT i.operation,
               i.changed,
               i.id,
               i.key,
               i.name,
               i.description,
               i.meta,
               i.tag,
               i.attribute,
               i.status,
               i.item_type_id,
               i.version,
               i.partition_id,
               i.created,
               i.updated,
               i.changed_by,
               it.key as item_type_key
        FROM ox_get_tree_content(root_item_key_param, label_param) s
               INNER JOIN item_change i
                          ON i.id = s.id::bigint
                            AND i.version = s.version::bigint
                            AND s.is_item = true
               INNER JOIN item_type_change it
                          ON it.id = i.item_type_id
                            AND it.version = i.version;
    END;
    $BODY$;

    ALTER FUNCTION ox_get_tree_items(character varying, character varying)
      OWNER TO onix;

    /*
      ox_get_tree_links(root_item_key_param, label_param): gets a list of all the links that are part
        of a tag tree for a specific parent item and a label.
     */
    CREATE OR REPLACE FUNCTION ox_get_tree_links(root_item_key_param character varying,
                                              label_param character varying)
      RETURNS TABLE
              (
                operation      character,
                changed        timestamp with time zone,
                id             bigint,
                key            character varying,
                link_type_id   integer,
                start_item_id  bigint,
                end_item_id    bigint,
                description    text,
                meta           jsonb,
                tag            text[],
                attribute      hstore,
                version        bigint,
                created        timestamp with time zone,
                updated        timestamp with time zone,
                changed_by     character varying,
                link_type_key  character varying,
                start_item_key character varying,
                end_item_key   character varying
              )
      LANGUAGE 'plpgsql'
      COST 100
      VOLATILE
      ROWS 1000
    AS
    $BODY$
    BEGIN
      RETURN QUERY
        SELECT l.*, lt.key as link_type_key, start_item.key as start_item_key, end_item.key as end_item_key
        FROM ox_get_tree_content(root_item_key_param, label_param) s
               INNER JOIN link_change l
                          ON l.id = s.id::bigint
                            AND l.version = s.version::bigint
                            AND s.is_item = false
               INNER JOIN link_type_change lt
                          ON lt.id = l.link_type_id
                            AND lt.version = l.version
               INNER JOIN item_change start_item
                          ON start_item.id = l.start_item_id
                            AND start_item.version = l.version
               INNER JOIN item_change end_item
                          ON end_item.id = l.end_item_id
                            AND end_item.version = l.version;

    END;
    $BODY$;

    ALTER FUNCTION ox_get_tree_links(character varying, character varying)
      OWNER TO onix;

    /*
      ox_get_item_tags(root_item_key_param): gets a list of tags for a specified items that
        is the parent for the tag tree.
     */
    CREATE OR REPLACE FUNCTION ox_get_item_tags(root_item_key_param character varying)
      RETURNS TABLE
              (
                id            integer,
                label         character varying,
                root_item_key character varying,
                name          character varying,
                description   text,
                item_data     hstore,
                link_data     hstore,
                version       bigint,
                created       timestamp with time zone,
                updated       timestamp with time zone,
                changed_by    character varying
              )
      LANGUAGE 'plpgsql'
      COST 100
      VOLATILE
      ROWS 1000
    AS
    $BODY$
    BEGIN
      RETURN QUERY
        SELECT *
        FROM tag s
        WHERE s.root_item_key = root_item_key_param;
    END;
    $BODY$;

    ALTER FUNCTION ox_get_item_tags(character varying)
      OWNER TO onix;

  END
  $$;