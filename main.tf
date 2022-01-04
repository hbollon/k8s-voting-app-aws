resource "aws_eks_cluster" "aws_eks" {
  name     = "eks_cluster_voting_app"
  role_arn = aws_iam_role.eks_cluster.arn

  vpc_config {
    subnet_ids = var.default_vpc_subnets
  }

  tags = {
    Name = "EKS_Voting_App_Stack"
  }
}

resource "aws_eks_node_group" "node" {
  cluster_name    = aws_eks_cluster.aws_eks.name
  node_group_name = "node-group-1"
  node_role_arn   = aws_iam_role.eks_nodes.arn
  subnet_ids      = var.default_vpc_subnets

  scaling_config {
    desired_size = 1
    max_size     = 1
    min_size     = 1
  }

  depends_on = [
    aws_iam_role_policy_attachment.AmazonEKSWorkerNodePolicy,
    aws_iam_role_policy_attachment.AmazonEKS_CNI_Policy,
    aws_iam_role_policy_attachment.AmazonEC2ContainerRegistryReadOnly,
  ]
}