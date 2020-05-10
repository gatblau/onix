---
id: concept_db
title: Database
---
Anything in the world we observe can be thought of being a "configuration" and can be recorded in one way or another.

Onix has a flexible database designed to record any configuration. This page explains how the Onix database is structured 
to record such configurations.

## Configurations

*Configuration* in the context of this document refers to the relative arrangement of parts or elements in an environment.

For example, a photograph is a point in time, visual representation of an arrangement of elements that appear in the picture. 

Any arrangement of elements could be thought of following a pattern that describes the form the arrangement takes. 
Such patterns are also called models. 

Conversely, models can be used to shape the form an arrangement should take by providing a convention and the basis for 
validation of the captured data. 

### Items & Links

In Onix, a part, or an element in an arrangement is called item (a.k.a. configuration item or CI). 

The arrangement itself is made of a combination of items and the links between them. This arrangement of items forms 
what graph theory calls a [graph structure](https://en.wikipedia.org/wiki/Graph_(discrete_mathematics)). 
The Onix logotype is a depiction of a 
graph structure as shown below: 
<img src="/onix/img/logo.png" width="200"/><img src="/onix/img/graph.png" width="200"/>

## Models

Arrangements of items and links in Onix follow a model. Models are fundamental to validate arrangements of linked items
in Onix.

To understand a model, consider the photograph example discussed earlier. Let us say that the elements in the picture are the 
members of a family. Such a family could be a stereotypical family, that is the family follows the the types of members 
and relationships between them that most people would expect.

For example:

- A family has a father
- A family has a mother
- A family has two children
- The family members are sitting on a sofa

Mother, father and children are human beings that play different roles. 

Mother, father and child can be thought of as types of humans (type in the sense of the roles they are playing).
Likewise, a sofa can be regarded as a type of furniture.

Then, the statement Mother is sitting on Sofa can be expressed by a link between the two items of different types.

However, the statement Sofa is sitting on Mother does not make sense. 

So there is a need for rules that constrain the way items are connected by links.

In Onix, the family model could be defined in the following way:

- Family is a *model*
- Mother is an *item type*
- Father is an *item type*
- A Child is an *item type*
- The relationship between a Mother and a Child is a *link type*
- The relationship between a Father and a Child is a *link type*
- The relationship between siblings is a *link type*
- The relationship between a Mother and a Sofa is a *link type*
- The relationship between a Father and a Sofa is a *link type*
- The constraint that a Mother can be sat on a Sofa is a *link rule*
- The fact that Mothers have a name is expressed as a attribute of the Mother *item type*
- etc...

Once the model is defined, Onix knows how information for different families can be recorded.

The actual items and links are depictions of different humans belonging in different families along with their relationships.

The following figure shows the entities in the Onix database and their relations:

![Onix Database](/onix/img/database.png)

### Item Types

Every time a new item is created in Onix, it needs the Item Type. The item type defines the characteristics of items of 
the type.

### Link Types

### Link Rules


