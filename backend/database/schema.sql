DROP TABLE IF EXISTS Notificacao CASCADE;
DROP TABLE IF EXISTS Pagamento CASCADE;
DROP TABLE IF EXISTS Reserva CASCADE;
DROP TABLE IF EXISTS AreaComum CASCADE;
DROP TABLE IF EXISTS Visitante CASCADE;
DROP TABLE IF EXISTS Comunicado CASCADE;
DROP TABLE IF EXISTS Visita CASCADE;
DROP TABLE IF EXISTS Usuario CASCADE;
DROP TABLE IF EXISTS Apartamento CASCADE;

CREATE TABLE Apartamento
(
  id uuid DEFAULT gen_random_uuid(),
  numero INTEGER NOT NULL,
  bloco VARCHAR(10) NOT NULL,
  CONSTRAINT PK_Apartamento PRIMARY KEY (id)
);

CREATE TABLE Usuario
(
  id uuid DEFAULT gen_random_uuid(),
  nome VARCHAR(100) NOT NULL,
  email VARCHAR(100) NOT NULL UNIQUE,
  senha VARCHAR(100) NOT NULL,
  telefone VARCHAR(15) NOT NULL,
  status VARCHAR(20) NOT NULL DEFAULT 'ATIVO',
  tipo VARCHAR(10) NOT NULL,
  apartamento_id uuid,
  responsavel BOOLEAN NOT NULL DEFAULT FALSE,
  CONSTRAINT PK_Usuario PRIMARY KEY (id),
  CONSTRAINT FK_Usuario_Apartamento FOREIGN KEY (apartamento_id) REFERENCES Apartamento (id),
  CONSTRAINT CK_Usuario_Status CHECK (status IN ('ATIVO', 'INATIVO', 'BLOQUEADO')),
  CONSTRAINT CK_Usuario_Tipo CHECK (tipo IN ('PORTEIRO', 'MORADOR', 'SINDICO')),
  CONSTRAINT CK_Usuario_ApartamentoConsistente CHECK (tipo = 'MORADOR' OR apartamento_id IS NULL),
  CONSTRAINT CK_Usuario_ResponsavelConsistente CHECK (responsavel = FALSE OR tipo = 'MORADOR')
);

CREATE TABLE Visitante
(
  id uuid DEFAULT gen_random_uuid(),
  nome VARCHAR(100) NOT NULL,
  cpf VARCHAR(11) NOT NULL UNIQUE,
  telefone VARCHAR(15) NOT NULL,
  foto VARCHAR(255),
  CONSTRAINT PK_Visitante PRIMARY KEY (id)
);

CREATE TABLE Visita
(
  id uuid DEFAULT gen_random_uuid(),
  dataEntrada TIMESTAMP NOT NULL,
  dataSaida TIMESTAMP,
  porteiro_id uuid NOT NULL,
  visitante_id uuid NOT NULL,
  morador_id uuid NOT NULL,
  CONSTRAINT PK_Visita PRIMARY KEY (id),
  CONSTRAINT FK_Visita_Porteiro FOREIGN KEY (porteiro_id) REFERENCES Usuario (id),
  CONSTRAINT FK_Visita_Visitante FOREIGN KEY (visitante_id) REFERENCES Visitante (id),
  CONSTRAINT FK_Visita_Morador FOREIGN KEY (morador_id) REFERENCES Usuario (id)
);

CREATE TABLE Comunicado
(
  id uuid DEFAULT gen_random_uuid(),
  titulo VARCHAR(100) NOT NULL,
  descricao TEXT NOT NULL,
  dataPublicacao TIMESTAMP NOT NULL,
  sindico_id uuid NOT NULL,
  CONSTRAINT PK_Comunicado PRIMARY KEY (id),
  CONSTRAINT FK_Comunicado_Sindico FOREIGN KEY (sindico_id) REFERENCES Usuario (id)
);

CREATE TABLE AreaComum
(
  id uuid DEFAULT gen_random_uuid(),
  nome VARCHAR(100) NOT NULL,
  descricao TEXT NOT NULL,
  capacidade INTEGER NOT NULL,
  CONSTRAINT PK_AreaComum PRIMARY KEY (id)
);

CREATE TABLE Reserva
(
  id uuid DEFAULT gen_random_uuid(),
  data DATE NOT NULL,
  horaInicio TIME NOT NULL,
  horaFim TIME NOT NULL,
  status VARCHAR(20) NOT NULL DEFAULT 'PENDENTE',
  morador_id uuid NOT NULL,
  areacomum_id uuid NOT NULL,
  sindico_id uuid,
  CONSTRAINT PK_Reserva PRIMARY KEY (id),
  CONSTRAINT FK_Reserva_Morador FOREIGN KEY (morador_id) REFERENCES Usuario (id),
  CONSTRAINT FK_Reserva_AreaComum FOREIGN KEY (areacomum_id) REFERENCES AreaComum (id),
  CONSTRAINT FK_Reserva_Sindico FOREIGN KEY (sindico_id) REFERENCES Usuario (id),
  CONSTRAINT CK_Reserva_Status CHECK (status IN ('PENDENTE', 'CONFIRMADA', 'CANCELADA'))
);

CREATE TABLE Pagamento
(
  id uuid DEFAULT gen_random_uuid(),
  valor DECIMAL(10, 2) NOT NULL,
  vencimento DATE NOT NULL,
  dataPagamento DATE,
  status VARCHAR(20) NOT NULL DEFAULT 'PENDENTE',
  morador_id uuid NOT NULL,
  CONSTRAINT PK_Pagamento PRIMARY KEY (id),
  CONSTRAINT FK_Pagamento_Morador FOREIGN KEY (morador_id) REFERENCES Usuario (id),
  CONSTRAINT CK_Pagamento_Status CHECK (status IN ('PENDENTE', 'PAGO', 'ATRASADO'))
);

CREATE TABLE Notificacao
(
  id uuid DEFAULT gen_random_uuid(),
  tipo VARCHAR(10) NOT NULL,
  destinatario_id uuid NOT NULL,
  mensagem TEXT NOT NULL,
  status VARCHAR(20) NOT NULL DEFAULT 'PENDENTE',
  dataEnvio TIMESTAMP,
  CONSTRAINT PK_Notificacao PRIMARY KEY (id),
  CONSTRAINT FK_Notificacao_Usuario FOREIGN KEY (destinatario_id) REFERENCES Usuario (id),
  CONSTRAINT CK_Notificacao_Tipo CHECK (tipo IN ('EMAIL', 'SMS')),
  CONSTRAINT CK_Notificacao_Status CHECK (status IN ('PENDENTE', 'ENVIADA', 'FALHA'))
);