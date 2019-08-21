## Face Recognition

1. You need to download the flowers classification samples from:
https://github.com/GoogleCloudPlatform/cloudml-samples/tree/master/tensorflow/standard/legacy/flowers

2. Edit /flowers/trainer/task.py file, update the default label_count to 2

```
parser.add_argument('--label_count', type=int, default=2)
```

3. Create run.sh under /flowers/ in cloud shell terminal:

```
$ chmod u+x run.sh
$ ./run.sh
```

4. To use a pretrained model, run:

```
gsutil cp -r gs://around-179500-mlengine/cc526__baggio/flowers_cc526__baggio_20180217_074840/training/model gs://around-250315-mlengine
```