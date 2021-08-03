# vmdiskfs command line tool

Tool used in artisan image conversion runtime. Utility has multiple usage:
- Simple conversion from any format to any format. Output format is customizable via environment variable
- Conversion and Deployment using config file, where future vm described. Example of that config file located in
  example folder
  
Utility can parse put notifications from almost any S3 type storage. Currently, only 2 types described:
- AWS S3
- MinIO

In the future enhancement planned to include Google Object Storage and Azure Object Storage put notification parsing.

Below are important variables required by utility:

|Environment Variable | Description                                   |
|---------------------|-----------------------------------------------|
|PIPELINE_HOME        | HOME directory                                |
|AWS_URL              | S3 URL                                        |
|AWS_ACCESS_KEY_ID    | S3 Access Key                                 |
|AWS_SECRET_ACCESS_KEY| S3 Secret Key                                 |
|AWS_DEST_BUCKET      | S3 Destination bucket name                    |
|AWS_PROVIDER         | S3 provider name, aws or minio                |
|AWS_USE_SSL          | S3 ssl usage, true or false                   |
|WEBHOOK_RECEIVER     | Artisan runner deploy vm webhook receiver url |
|OUTPUT_FORMAT        | output format: qcow2 or vhd or other          |

Current directory from where vmdiskfs utility is running must have scripts, source and output 
directories.

Note: Convert requires qemu-img tool to be installed

- Building vmdifkfs tool for Linux and MacOS
```shell
#For MacOS
art run build-mac
#For Linux
art run build-linux
```
- Archiving to zip or tar.gz
```shell
#Use on MacOS
art run zip
#Use on any OS except Windows
art run tar
```
- Make vmdiskfs package
```shell
art build -t vmdiskfs
```