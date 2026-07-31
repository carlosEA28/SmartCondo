resource "aws_db_subnet_group" "this" {
  name       = "${var.project_name}-${var.environment}"
  subnet_ids = var.private_subnet_ids

  tags = {
    Name = "${var.project_name}-${var.environment}"
  }
}

resource "aws_rds_cluster" "this" {
  cluster_identifier = "${var.project_name}-${var.environment}"
  engine             = "aurora-postgresql"
  engine_version     = "15.4"
  database_name      = var.db_name
  master_username    = var.db_username
  master_password    = var.db_password

  db_subnet_group_name   = aws_db_subnet_group.this.name
  vpc_security_group_ids = [var.rds_security_group_id]

  storage_encrypted   = true
  skip_final_snapshot = var.environment != "prod"

  tags = {
    Name = "${var.project_name}-${var.environment}"
  }
}

resource "aws_rds_cluster_instance" "this" {
  count              = var.environment == "prod" ? 2 : 1
  identifier         = "${var.project_name}-${var.environment}-${count.index}"
  cluster_identifier = aws_rds_cluster.this.id
  instance_class     = var.instance_class
  engine             = aws_rds_cluster.this.engine
  engine_version     = aws_rds_cluster.this.engine_version

  tags = {
    Name = "${var.project_name}-${var.environment}-${count.index}"
  }
}
