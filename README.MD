# Telegram to VCF

A small utility to extract your Telegram contacts as VCF file.

## Usage

At first build this application or use the precompiled binaries on releases page to download the program. Then you have
to set some environment variables through your terminal:

```bash
export APP_ID=1234
export APP_HASH=abcd
export PHONE="+989123456789"
export DOWNLOAD_PHOTOS=true # Optional; There is no need to have this if you don't want the pictures
```

Then simply run the program and wait for it to complete. The output file will be put in the working directory of your
program in a file named `contacts.vcf`.

### What are APP_ID and APP_HASH?

You can obtain these values from [here](https://my.telegram.org/).