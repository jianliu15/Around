declare -r PROJECT=$(gcloud config list project --format "value(core.project)")
declare -r JOB_ID="face_${USER}_$(date +%Y%m%d_%H%M%S)"
declare -r BUCKET="gs://${PROJECT}-mlengine"

declare -r GCS_PATH="${BUCKET}/${USER}/${JOB_ID}"
declare -r DICT_FILE=gs://face_179500/dict.txt

declare -r MODEL_NAME=face
declare -r VERSION_NAME=v1

python trainer/preprocess.py \
  --input_dict "$DICT_FILE" \
  --input_path "gs://face_179500/eval_set.csv" \
  --output_path "${GCS_PATH}/preproc/eval" \
  --cloud

python trainer/preprocess.py \
  --input_dict "$DICT_FILE" \
  --input_path "gs://face_179500/train_set.csv" \
  --output_path "${GCS_PATH}/preproc/train" \
  --cloud

gcloud ml-engine jobs submit training "$JOB_ID" \
  --stream-logs \
  --module-name trainer.task \
  --package-path trainer \
  --staging-bucket "$BUCKET" \
  --region us-central1 \
  --runtime-version=1.0 \
  -- \
  --max-steps 300 \
  --output_path "${GCS_PATH}/training" \
  --eval_data_paths "${GCS_PATH}/preproc/eval*" \
  --train_data_paths "${GCS_PATH}/preproc/train*"

# Remove the model and its version
# Make sure no error is reported if model does not exist
gcloud ml-engine versions delete $VERSION_NAME --model=$MODEL_NAME -q --verbosity none
gcloud ml-engine models delete $MODEL_NAME -q --verbosity none

# Tell CloudML about a new type of model coming.  Think of a "model" here as
# a namespace for deployed Tensorflow graphs.
gcloud ml-engine models create "$MODEL_NAME" \
  --regions us-central1

# Each unique Tensorflow graph--with all the information it needs to execute--
# corresponds to a "version".  Creating a version actually deploys our
# Tensorflow graph to a Cloud instance, and gets is ready to serve (predict).
gcloud ml-engine versions create "$VERSION_NAME" \
  --model "$MODEL_NAME" \
  --origin "${GCS_PATH}/training/model" \
  --runtime-version=1.0
