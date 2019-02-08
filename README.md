# s3s
`s3s` is a tiny CLI that helps you upload local file to your private s3 bucket, generating a presign URL for this file followed by URL shortening with `bit.ly`.

## Features
- [x] simple CLI interface
- [x] support both AWS global regions as well as regions in China(`Beijing` and `Ningxia`).
- [x] `bit.ly` as the default URL shorterning provider
- [x] `t.cn` as the default provider for AWS China regions upload(`cn-north-1` and `cn-northwest-2`) [#1](https://github.com/pahud/s3s/issues/1)
- [x] single executable binary. No dependency required.
 


## Build

```
$ make build
$ cp ./s3s ~/bin/
$ s3s
usage: s3s <bucket> <filename>
```

add `export BITLY_TOKEN='<YOUR_BITLY_TOKEN>'` in your `~/.bash_profile`                                                                                                                         


## Usage

```
AWS_PROFILE=<YOUR_PROFILE_NAME> AWS_DEFAULT_REGION=<TARGET_AWS_REGION> s3s <YOUR_S3_BCKET> FILE  
```


## Example #1: through AWS Global regions
```
$ AWS_DEFAULT_REGION=ap-northeast-1 s3s pahud-tmp-nrt  ../lambda-layer-awscli/func-bundle.zip                                                             
Uploading file to S3...
Successfully uploaded ../lambda-layer-awscli/layer.zip to https://pahud-tmp-ap-northeast-1.s3.ap-northeast-1.amazonaws.com/layer.zip
Original URL: https://pahud-tmp-ap-northeast-1.s3.ap-northeast-1.amazonaws.com/layer.zip?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=AKIAJJFINOFXR3LTLIFQ%2F20190208%2Fap-northeast-1%2Fs3%2Faws4_request&X-Amz-Date=20190208T170601Z&X-Amz-Expires=604800&X-Amz-SignedHeaders=host&X-Amz-Signature=99034e25d36fa0200e336bfd17727b9ca36fe1a2452cefb65e3022eaea60cd3f
bitly URL: http://bit.ly/2DopVeo      
```

## Example #2: through AWS China Beijing or Ningxia region

```
$ AWS_PROFILE=cn AWS_DEFAULT_REGION=cn-northwest-1 s3s pahud-tmp-cn-northwest-1  ../lambda-layer-awscli/func-bundle.zip                                   
Uploading file to S3...
Successfully uploaded ../lambda-layer-awscli/func-bundle.zip to https://pahud-tmp-cn-northwest-1.s3.cn-northwest-1.amazonaws.com.cn/func-bundle.zip
Original URL: https://pahud-tmp-cn-northwest-1.s3.cn-northwest-1.amazonaws.com.cn/func-bundle.zip?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=AKIAOEVQOGI2VHQNTNHA%2F20190208%2Fcn-northwest-1%2Fs3%2Faws4_request&X-Amz-Date=20190208T170640Z&X-Amz-Expires=604800&X-Amz-SignedHeaders=host&X-Amz-Signature=b1d276e7fb4267b448c87977fae1b32f60be2835781727eed133bd1e2de96456
t.cn URL: http://t.cn/EcXukK5        
```

You may also set an alias in your `~/.bash_profile` like this

```
alias s3sjp="AWS_PROFILE=default AWS_DEFAULT_REGION=ap-northeast-1 $HOME/bin/s3s pahud-tmp-ap-northeast-1 $@"                                                                                             
alias s3scn="AWS_PROFILE=cn AWS_DEFAULT_REGION=cn-northwest-1 $HOME/bin/s3s pahud-tmp-cn-northwest-1 $@"                                                                                             
```

And simply 

```
$ s3sjp FILE
or
$ s3scn FILE
```

to upload your local file to your private S3 bucket and get the `bitly` URL immediately.

e.g.

upload to `Tokyo region(ap-northeast-1)`
```
$ s3sjp ../lambda-layer-awscli/func-bundle.zip                                                                                                            
Uploading file to S3...
Successfully uploaded ../lambda-layer-awscli/func-bundle.zip to https://pahud-tmp-ap-northeast-1.s3.ap-northeast-1.amazonaws.com/func-bundle.zip
Original URL: https://pahud-tmp-ap-northeast-1.s3.ap-northeast-1.amazonaws.com/func-bundle.zip?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=AKIAJJFINOFXR3LTLIFQ%2F20190208%2Fap-northeast-1%2Fs3%2Faws4_request&X-Amz-Date=20190208T170707Z&X-Amz-Expires=604800&X-Amz-SignedHeaders=host&X-Amz-Signature=7c3e9d08b97f4958fdf474b1cbdefe3adceb93ade8557002b6ef6724e2126130
bitly URL: http://bit.ly/2DivtHq        
```

upload to `Ningxia region(cn-northwest-1)`

```
$ s3scn ../lambda-layer-awscli/func-bundle.zip                                                                                                            
Uploading file to S3...
Successfully uploaded ../lambda-layer-awscli/func-bundle.zip to https://pahud-tmp-cn-northwest-1.s3.cn-northwest-1.amazonaws.com.cn/func-bundle.zip
Original URL: https://pahud-tmp-cn-northwest-1.s3.cn-northwest-1.amazonaws.com.cn/func-bundle.zip?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=AKIAOEVQOGI2VHQNTNHA%2F20190208%2Fcn-northwest-1%2Fs3%2Faws4_request&X-Amz-Date=20190208T170722Z&X-Amz-Expires=604800&X-Amz-SignedHeaders=host&X-Amz-Signature=326138a370657953b0e1fcd5f8243a94bc973a80418bda8f8a5fa035b37a2d72
t.cn URL: http://t.cn/EcX3Ah1                                                                                   
```