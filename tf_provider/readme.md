# Onix Provider <img src="../docs/pics/ox.png" width="200" height="200" align="right">

The provider documentation can be found [here](docs/index.md).

## Getting Started 

- Build the provider binary 
    - Requires GNU make
    - Requires Golang (1.13.x)
    - Requires Terraform (0.12.x)
    - run "make package"
    - The provider files are located within the build folder
    
- Install the provider
  - The provider must be manually installed, since "terraform init" cannot automatically 
  download it. 
  
  - Place the binary for the target operating system in the user plugins directory. 
  
The user plugins directory is in one of the following locations, depending on the host operating system:

|Operating system|	User plugins directory|
|---|---|
|Windows	| %APPDATA%\terraform.d\plugins|
|All other systems|	~/.terraform.d/plugins|

Alternatively, the provider can also be placed under the local folder where Terraform is run from. 

For example: **working_directory/.terraform.d/plugins**

Once the plugin is placed in the correct location, the [terraform init](https://www.terraform.io/docs/commands/init.html) command has to be run to initialise the working directory ready for use.


