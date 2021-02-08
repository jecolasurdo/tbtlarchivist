
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "tbtlarchivist_rg" {
  name     = "tbtlarchivist_rg"
  location = "eastus"

  tags = {
    environment = "tbtlarchivist"
  }
}

resource "azurerm_virtual_network" "tbtlarchivist_vn" {
  name                = "tbtlarchivist_vn"
  address_space       = ["10.0.0.0/16"]
  location            = "eastus"
  resource_group_name = azurerm_resource_group.tbtlarchivist_rg.name

  tags = {
    environment = "tbtlarchivist"
  }
}

resource "azurerm_subnet" "tbtlarchivist_sn" {
  name                 = "tbtlarchivist_sn"
  resource_group_name  = azurerm_resource_group.tbtlarchivist_rg.name
  virtual_network_name = azurerm_virtual_network.tbtlarchivist_vn.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_public_ip" "tbtlarchivist_public_ip" {
  name                = "public_ip"
  location            = "eastus"
  resource_group_name = azurerm_resource_group.tbtlarchivist_rg.name
  allocation_method   = "Dynamic"

  tags = {
    environment = "tbtlarchivist"
  }
}

resource "azurerm_network_security_group" "tbtlarchivist_sg" {
  name                = "tbtlarchivist_sg"
  location            = "eastus"
  resource_group_name = azurerm_resource_group.tbtlarchivist_rg.name

  security_rule {
    name                       = "SSH"
    priority                   = 1001
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "22"
    source_address_prefix      = "*"
    destination_address_prefix = "*"
  }

  tags = {
    environment = "tbtlarchivist"
  }
}

resource "azurerm_network_interface" "tbtlarchivist_nic" {
  name                = "tbtlarchivist_nic"
  location            = "eastus"
  resource_group_name = azurerm_resource_group.tbtlarchivist_rg.name

  ip_configuration {
    name                          = "tbtlarchivist_ip_config"
    subnet_id                     = azurerm_subnet.tbtlarchivist_sn.id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.tbtlarchivist_public_ip.id
  }

  tags = {
    environment = "tbtlarchivist"
  }
}

resource "azurerm_network_interface_security_group_association" "tbtlarchivist_sg_association" {
  network_interface_id      = azurerm_network_interface.tbtlarchivist_nic.id
  network_security_group_id = azurerm_network_security_group.tbtlarchivist_sg.id
}

resource "random_id" "randomId" {
  keepers = {
    resource_group = azurerm_resource_group.tbtlarchivist_rg.name
  }

  byte_length = 8
}

resource "azurerm_storage_account" "tbtlarchivist_sa" {
  name                     = "diag${random_id.randomId.hex}"
  resource_group_name      = azurerm_resource_group.tbtlarchivist_rg.name
  location                 = "eastus"
  account_replication_type = "LRS"
  account_tier             = "Standard"

  tags = {
    environment = "tbtlarchivist"
  }
}


resource "tls_private_key" "test_ssh" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

output "tls_private_key" { value = tls_private_key.test_ssh.private_key_pem }

resource "azurerm_linux_virtual_machine" "tbtlarchivist_vm" {
  name                  = "tbtlarchivist_vm"
  location              = "eastus"
  resource_group_name   = azurerm_resource_group.tbtlarchivist_rg.name
  network_interface_ids = [azurerm_network_interface.tbtlarchivist_nic.id]
  size                  = "Standard_B1s"

  os_disk {
    name                 = "tbtlarchivist_os_disk"
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "OpenLogic"
    offer     = "CentOS"
    sku       = "8_2"
    version   = "latest"
  }

  computer_name                   = "tbtlarchivist1"
  admin_username                  = "joe"
  disable_password_authentication = true

  admin_ssh_key {
    username   = "joe"
    public_key = tls_private_key.test_ssh.public_key_openssh
  }

  boot_diagnostics {
    storage_account_uri = azurerm_storage_account.tbtlarchivist_sa.primary_blob_endpoint
  }

  tags = {
    environment = "tbtlarchivist"
  }
}
