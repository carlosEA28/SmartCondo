output "endpoint" {
  value     = aws_rds_cluster.this.endpoint
  sensitive = true
}

output "port" {
  value = aws_rds_cluster.this.port
}

output "database_name" {
  value = aws_rds_cluster.this.database_name
}
