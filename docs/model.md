# CMDB Model

<a name="toc"></a>
## Table of Contents [(index)](./../readme.md)

- [Semantic Model](#semantic-model)
- [Logical Model](#relational-model)


<a name="semantic-model"></a>
## Semantic Model [(up)](#toc)

The following figure shows the [semantic model](https://en.wikipedia.org/wiki/Semantic_data_model) for the CMDB:
 
![Semantic Data Model](./pics/semantic_model.png "Onix Semantic Data Model")

- **Items** store configuration information and can be associated to other items using **Links**.
- **Items** are of a specified **Item Type**.
- **Links** connect **Items** creating associations.
- **Links** are of a specified **Link Type**.
- **Link Rules** apply to particular **Links** and restrict what **Item Types** the **Link** can connect.
- **Models** are collections of **Item Types** and **Link Types**.

<a name="relational-model"></a>
## Logical Model [(up)](#toc)

The following picture shows the Onix logical data model:

![Logical Data Model](./pics/logical_model.png "Onix Relational Data Model")



