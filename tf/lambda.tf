resource "aws_iam_role" "role_lambda_backend" {
  name               = "ctrlv-lambda-role-${local.unique_id}"
  description        = "execution role for ctrlv backend"
  assume_role_policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Principal": {
                "Service": "lambda.amazonaws.com"
            },
            "Action": "sts:AssumeRole"
        }
    ]
}
EOF
}

resource "aws_iam_policy" "policy_lambda_backend" {
  name        = "ctrlv-lambda-policy-${local.unique_id}"
  description = "policy for backend, granting dynamodb and logs"
  policy      = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "DynamoDBTableAccess",
            "Effect": "Allow",
            "Action": [
                "dynamodb:Scan",
                "dynamodb:Query",
                "dynamodb:GetItem",
                "dynamodb:PutItem",
                "dynamodb:DeleteItem",
                "dynamodb:UpdateItem"
            ],
            "Resource": "${aws_dynamodb_table.ctrlv_db.arn}"
        },
        {
            "Effect": "Allow",
            "Action": [
                "logs:CreateLogGroup",
                "logs:CreateLogStream",
                "logs:PutLogEvents"
            ],
            "Resource": "*"
        }
    ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "attach_backend_policy" {
  role       = aws_iam_role.role_lambda_backend.name
  policy_arn = aws_iam_policy.policy_lambda_backend.arn
}

resource "aws_lambda_function" "lambda_ctrlv_backend" {
  filename      = "../bin/main.zip"
  function_name = "ctrlv-lambda-${local.unique_id}"
  description   = "ctrlv server"
  role          = aws_iam_role.role_lambda_backend.arn
  handler       = "main"
  runtime       = "go1.x"
  memory_size   = 128 # 128 MB to 10,240 MB
  timeout       = 5   # 5s to 900s
}
