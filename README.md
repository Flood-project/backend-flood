# 📦 Backend - Catálogo de Produtos Robustec

API RESTful desenvolvida em **Go** para gerenciamento do catálogo de produtos da **Robustec LTDA**, especializada em pés de apoio industriais.

---

## 📋 Descrição do Projeto

O **Backend do Catálogo de Produtos Robustec** é uma API moderna com autenticação, auditoria, filtros avançados e armazenamento de arquivos.  
Ele oferece:

- 🔐 Autenticação JWT (access/refresh)
- 📦 CRUD completo de produtos e componentes
- 🗃️ Upload de arquivos no MinIO
- 🧱 Arquitetura limpa e escalável
- 🐳 Deploy com Docker/Docker Compose  
- 📝 Auditoria automática via PostgreSQL

---

## 🏗️ Arquitetura e Tecnologias

### 🐹 Backend (Go)
- **Go**
- **Chi Router**
- **SQLx**
- **JWT**
- **Bcrypt**
- **CORS**

### 🐳 Infraestrutura
- **Docker**
- **Docker Compose**
- **PostgreSQL 14**
- **MinIO (S3 compatível)**

### 🗃️ Banco e Armazenamento
- PostgreSQL para dados estruturados  
- MinIO para arquivos e imagens  

---

## ⚙️ Requisitos para Rodar Localmente

### Pré-requisitos
- Docker + Docker Compose *(recomendado)*
- Go 1.21+ *(opcional para desenvolvimento local)*
- Git

### Configuração `.env`

Crie um arquivo `.env`:

```env

    # Banco de Dados
    DB_HOST=postgres
    DATABASE_URL=postgresql://
    DB_PORT=5432
    DB_USER=postgres
    DB_PASSWORD=postgres
    DB_NAME=robustec_db
    DB_SSLMODE=require

    # JWT
    SECRET_KEY=sua_chave_secreta_jwt_aqui

    # MinIO
    MINIO_ENDPOINT=minio:9000
    MINIO_ACCESS_KEY=minioadmin
    MINIO_SECRET_KEY=minioadmin
```

# CORS
ALLOWED_ORIGINS=http://localhost:5173,http://localhost:3000


## ▶️ Como Rodar o Projeto

### 🚀 Usando Docker (recomendado)

Para subir toda a stack (API + PostgreSQL + MinIO):

```bash
docker-compose up --build
```

### 📡 Endpoints Após Subir o Docker

API Backend:
http://localhost:8080

MinIO (console administrativo):
http://localhost:9000

Usuário: minioadmin
Senha: minioadmin

### 📡 Endpoints Principais da API

Abaixo está a lista dos endpoints mais utilizados após subir o Docker 👇

---

## 🔐 Autenticação

### **POST /login**  
Realiza o login e retorna tokens JWT (access e refresh).


---

## 📦 Produtos

### **POST /products**  
Cria um novo produto.

### **GET /products**  
Lista produtos com suporte a filtros, paginação e componentes.

### **GET /products/{id}**  
Obtém um produto pelo ID.

### **PUT /products/{id}**  
Atualiza os dados de um produto.

### **DELETE /products/{id}**  
Remove um produto pelo ID.

### 🔍 Exemplo com filtros:

GET /products?eq[tipobucha]=1&eq[tipoacionamento]=2&page=1&limit=10


---

## 🗂️ Upload de Arquivos

### **POST /files/{ID_do_produto}**  
Realiza upload de imagens/arquivos para o MinIO.


## 📝 Auditoria

### **GET /logs**  
Lista os logs gerados automaticamente pelo sistema.

---

### 🌐 **Interface Pública (Catálogo)**
- Visualização de produtos ativos com imagens, especificações técnicas e descrições
- Sistema de filtros avançados (por tipo de bucha, acionamento e base)
- Busca por código de produto
- Design responsivo para mobile, tablet e desktop
- Carrossel de imagens para cada produto
- Detalhes completos dos produtos em modal
- Integração com WhatsApp para contato direto

### 🔐 **Interface Administrativa**
- CRUD completo de produtos (Criar, Visualizar, Atualizar, Deletar)
- Upload e gerenciamento de múltiplas imagens por produto
- Sistema de logs detalhado de todas as operações
- Gerenciamento de usuários e permissões
- Controle de componentes (buchas, acionamentos, bases)
- Autenticação JWT com níveis de acesso

### 👤 **Interface Comercial**
- CRUD completo de usuários (Criar, Visualizar, Atualizar e Deletar)
- Gerenciamento de usuários e permissões
- Visualização da tela de cliente
- Autenticação JWT com níveis de acesso
---

## 📞 Suporte

-Documentação e testes -> Jailopesoutlook@gmail.com
-Backend e testes -> apspolti@gmail.com
-Frontend e testes -> eduardoosartori@gmail.com

---

## 🗺️ Roadmap

Vários tópicos ainda restam para serem tratados, como:

- Paginação de Usuários, Componentes e Auditoria;
- Melhorias gerais de UX e validações no front e back;
- Pesquisa de produtos por nome personalizado;
- Filtros de auditoria e pesquisa por usuário/produto alterado e responsável pela alteração;
- Ajustar armazenamento das imagens e edição das imagens. Quando usuário for editar, que mostre a imagem atual e permita subsituí-la;
- Adicionar histórico de eventos nos usuários e produtos (logs diretamente em cada produto ou usuário separadamente, sem necessidade de ver todos os logs);
- Adicionar tipos de categorias de produtos diferentes;
- Permitir criar tipos de usuário e permissões personalizadas;
- Melhorar retornos de auditoria, mostrando exatamente o que havia antes e o que restou depois de alterações;
- Otimizar código visando velocidade;
- Redefinição de senha de usuários via e-mail
- E muito mais....

---

## 👥 Autores

-Documentação e testes -> Jailopesoutlook@gmail.com e 
danielsoranco@cesurg.com
-Backend e testes -> apspolti@gmail.com
-Frontend e testes -> eduardoosartori@gmail.com

---

## 📄 Licença

**Proprietária** - Todos os direitos reservados © Robustec Indústria e Comércio Ltda

Este projeto é de propriedade exclusiva da Robustec LTDA. Nenhuma parte deste software pode ser reproduzida, distribuída ou transmitida de qualquer forma ou por qualquer meio sem a permissão prévia por escrito da Robustec LTDA.

---

## 📊 Status do Projeto

🟢 **Em Desenvolvimento Ativo**

- ✅ Primeira versão do catálogo concluída
- ✅ Sistema administrativo funcional
- ✅ Autenticação e autorização implementadas
- ✅ Sistema de logs operacional
- 🔄 Melhorias contínuas e novas funcionalidades em andamento

---

## 🌟 Funcionalidades Principais

### ✨ Catálogo Público
- [x] Listagem de produtos com paginação
- [x] Filtros por componentes (bucha, acionamento, base)
- [x] Busca por código de produto
- [x] Visualização de detalhes com múltiplas imagens
- [x] Design responsivo (mobile-first)
- [x] Integração com WhatsApp

### 🔐 Área Administrativa
- [x] CRUD completo de produtos - Criar, Alterar, Excluir e Visualizar
- [x] Upload de múltiplas imagens
- [x] Sistema de logs detalhado
- [x] Gerenciamento de usuários
- [x] Controle de componentes - Buchas, Acionamentos e Bases
- [x] Autenticação JWT

---

**Desenvolvido com ❤️ para a equipe Robustec**