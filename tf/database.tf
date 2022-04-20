resource "aws_dynamodb_table" "ctrlv_db" {
  name           = local.dynamodb_table
  billing_mode   = "PROVISIONED"
  read_capacity  = 1
  write_capacity = 1
  hash_key       = "Id"

  attribute {
    name = "Id"
    type = "S"
  }

  attribute {
    name = "Alias"
    type = "S"
  }

  global_secondary_index {
    name            = local.dynamodb_table_gsi
    hash_key        = "Alias"
    write_capacity  = 1
    read_capacity   = 1
    projection_type = "KEYS_ONLY"
  }

  tags = {
    app = "ctrlv"
  }
}

# ------------ Auto scaling for table
resource "aws_appautoscaling_target" "dynamodb_table_read_target" {
  min_capacity       = 1
  max_capacity       = 10
  resource_id        = "table/${local.dynamodb_table}"
  scalable_dimension = "dynamodb:table:ReadCapacityUnits"
  service_namespace  = "dynamodb"
  depends_on         = [aws_dynamodb_table.ctrlv_db]
}

resource "aws_appautoscaling_policy" "dynamodb_table_read_policy" {
  name               = "DynamoDBReadCapacityUtilization:${aws_appautoscaling_target.dynamodb_table_read_target.resource_id}"
  policy_type        = "TargetTrackingScaling"
  resource_id        = aws_appautoscaling_target.dynamodb_table_read_target.resource_id
  scalable_dimension = aws_appautoscaling_target.dynamodb_table_read_target.scalable_dimension
  service_namespace  = aws_appautoscaling_target.dynamodb_table_read_target.service_namespace
  depends_on         = [aws_dynamodb_table.ctrlv_db]

  target_tracking_scaling_policy_configuration {
    predefined_metric_specification {
      predefined_metric_type = "DynamoDBReadCapacityUtilization"
    }

    target_value = 70
  }
}

resource "aws_appautoscaling_target" "dynamodb_table_write_target" {
  min_capacity       = 1
  max_capacity       = 10
  resource_id        = "table/${local.dynamodb_table}"
  scalable_dimension = "dynamodb:table:WriteCapacityUnits"
  service_namespace  = "dynamodb"
  depends_on         = [aws_dynamodb_table.ctrlv_db]
}

resource "aws_appautoscaling_policy" "dynamodb_table_write_policy" {
  name               = "DynamoDBWriteCapacityUtilization:${aws_appautoscaling_target.dynamodb_table_write_target.resource_id}"
  policy_type        = "TargetTrackingScaling"
  resource_id        = aws_appautoscaling_target.dynamodb_table_write_target.resource_id
  scalable_dimension = aws_appautoscaling_target.dynamodb_table_write_target.scalable_dimension
  service_namespace  = aws_appautoscaling_target.dynamodb_table_write_target.service_namespace
  depends_on         = [aws_dynamodb_table.ctrlv_db]

  target_tracking_scaling_policy_configuration {
    predefined_metric_specification {
      predefined_metric_type = "DynamoDBWriteCapacityUtilization"
    }

    target_value = 70
  }
}

# ------------ Auto scaling for global secondary index

resource "aws_appautoscaling_target" "dynamodb_gsi_read_target" {
  max_capacity       = 10
  min_capacity       = 1
  resource_id        = "table/${local.dynamodb_table}/index/${local.dynamodb_table_gsi}"
  scalable_dimension = "dynamodb:index:ReadCapacityUnits"
  service_namespace  = "dynamodb"
  depends_on         = [aws_dynamodb_table.ctrlv_db]
}

resource "aws_appautoscaling_policy" "dynamodb_gsi_read_policy" {
  name               = "DynamoDBReadCapacityUtilization:${aws_appautoscaling_target.dynamodb_gsi_read_target.resource_id}"
  policy_type        = "TargetTrackingScaling"
  resource_id        = aws_appautoscaling_target.dynamodb_gsi_read_target.resource_id
  scalable_dimension = aws_appautoscaling_target.dynamodb_gsi_read_target.scalable_dimension
  service_namespace  = aws_appautoscaling_target.dynamodb_gsi_read_target.service_namespace
  depends_on         = [aws_dynamodb_table.ctrlv_db]

  target_tracking_scaling_policy_configuration {
    predefined_metric_specification {
      predefined_metric_type = "DynamoDBReadCapacityUtilization"
    }

    target_value = 70
  }
}

resource "aws_appautoscaling_target" "dynamodb_gsi_write_target" {
  max_capacity       = 10
  min_capacity       = 1
  resource_id        = "table/${local.dynamodb_table}/index/${local.dynamodb_table_gsi}"
  scalable_dimension = "dynamodb:index:WriteCapacityUnits"
  service_namespace  = "dynamodb"
  depends_on         = [aws_dynamodb_table.ctrlv_db]
}

resource "aws_appautoscaling_policy" "dynamodb_gsi_write_policy" {
  name               = "DynamoDBWriteCapacityUtilization:${aws_appautoscaling_target.dynamodb_gsi_write_target.resource_id}"
  policy_type        = "TargetTrackingScaling"
  resource_id        = aws_appautoscaling_target.dynamodb_gsi_write_target.resource_id
  scalable_dimension = aws_appautoscaling_target.dynamodb_gsi_write_target.scalable_dimension
  service_namespace  = aws_appautoscaling_target.dynamodb_gsi_write_target.service_namespace
  depends_on         = [aws_dynamodb_table.ctrlv_db]

  target_tracking_scaling_policy_configuration {
    predefined_metric_specification {
      predefined_metric_type = "DynamoDBWriteCapacityUtilization"
    }

    target_value = 70
  }
}
