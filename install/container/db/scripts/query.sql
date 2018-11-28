-- findNodesByTypeAndTag
-- find all nodes of a particular type that have a particular tag
SELECT *
FROM node
WHERE tag LIKE '%Prod%'

-- findLinkedNodesByTypeAndTag
-- find all nodes of a particular type and tag which are linked to the specified node
SELECT endNode.*
FROM link
	INNER JOIN node endNode
		ON endNode.id = link.end_node_id
		AND endNode.tag LIKE '%%' -- filtering tags
		AND endNode.node_type_id = 2 -- end node type
	INNER JOIN node startNode
		ON startNode.id = link.start_node_id
		AND startNode.id = 8 -- the id of the start node used to find linked nodes
ORDER BY startNode.id ASC

-- findAnyLinkedNodes
-- find all nodes linked to the specified node
SELECT endNode.*
FROM link
	INNER JOIN node endNode
		ON endNode.id = link.end_node_id
		AND endNode.tag LIKE '%%' -- filtering tags
	INNER JOIN node startNode
		ON startNode.id = link.start_node_id
		AND startNode.id = 8 -- the id of the start node used to find linked nodes
ORDER BY startNode.id ASC

-- findNodeTypesLinkedToNode
-- find all node types which are linked to a particular node
SELECT DISTINCT t.*
FROM link
	INNER JOIN node endNode
		ON endNode.id = link.end_node_id
	INNER JOIN node startNode
		ON startNode.id = link.start_node_id
		AND startNode.id = 8 -- the id of the start node used to find linked nodes
	INNER JOIN node_type t
		ON endNode.node_type_id = t.id
ORDER BY t.id ASC

-- findAllLinksFromNode
-- find all links coming off a particular node
SELECT *
FROM link
WHERE start_node_id = 8

-- findAllLinksToNode
-- find all links coming into a particular node
SELECT *
FROM link
WHERE end_node_id = 11

-- tags for searching
INSERT INTO public.item(name, description, status, item_type_id, meta, version, created, updated, tag, key)
VALUES ('n', 'd', 0, 1, '{}', 0, '2018-11-28', '2018-11-28', ARRAY['test', 'cmdb'], 'key1');

SELECT *
FROM item
WHERE tag @> ARRAY['cmdb2', 'test']
