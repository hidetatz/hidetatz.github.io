Lambda Function to convert PDF to Image by Node.js---2018-04-12 21:28:36

I've created AWS Lambda function to convert pdf to image.

Repository is [here](https://github.com/ygnmhdtt/lmd_pdf2img).  
I need this function for work, but I couldn't find any useful libraries or examples.

### How to use
First, you need to install dependencies.

```
$ npm install gm aws-sdk util
```

Now, I haven't prepared any deployment command.  
You need to do it manually like this.

```
$ zip -r function.zip *
```

And, you need to upload zip at lambda console.  
Of course, you must set function trigger `S3 Upload` .

### Source code

All in `index.js` .

```
const srcBucket = event.Records[0].s3.bucket.name;
const srcKey = decodeURIComponent(event.Records[0].s3.object.key.replace(/\+/g, " "));

const dstBucket = srcBucket;
const dstKey = srcKey.replace('.pdf', '.png');
```

Now, destination bucket and key is same as source (extension will be `png` ).  
If you want to configure it, please write your own code.

Other configurations are written in README.  
