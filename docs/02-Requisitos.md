# Requisitos — SmartCondo

---

## Problemas que o sistema resolve

| Problema | O que acontece |
|---|---|
| Não tem lembrete de pagamento | Muita gente não paga a tempo |
| Reserva sem controle | Dois moradores marcam o mesmo horário |
| Visitante demora na portaria | Processo manual para liberar |
| Falta de comunicação | Morador não sabe de obras e manutenções |

---

## Quem usa o sistema

| Pessoa | O que faz |
|---|---|
| **Morador** | Reserva áreas comuns, cadastra visitantes, vê pagamentos e comunicados |
| **Porteiro** | Vê visitantes cadastrados e registra entrada/saída |
| **Síndico** | Posta comunicados, gerencia reservas |
| **Administrador** | Gerencia usuários e configura o sistema |

---

## Regras de negócio

| Regra |
|---|
| Mandar e-mail e SMS **7 dias antes** do pagamento vencer |
| Se não pagar, mandar aviso **no dia seguinte** |
| Morador que não pagou fica com status **"Inadimplente"** |
| Não pode reservar se já tiver reserva no **mesmo lugar, dia e horário** |
| Morador inadimplente **não pode** reservar |
| Só o **responsável pelo apartamento** pode reservar |
| Só o **responsável pelo apartamento** pode cadastrar visitante |
| O porteiro vê os visitantes que já foram cadastrados |
| O porteiro registra quando o visitante entra e sai |
| O síndico pode **aprovar, editar ou cancelar** reservas |
| O síndico posta comunicados sobre obras e manutenções |
| Cada pessoa acessa só o que precisa (morador, porteiro, síndico, admin) |

---

## Requisitos funcionais

| O que o sistema deve fazer |
|---|
| Ter um **dashboard** (painel) que funcione no computador e celular |
| Permitir que os moradores **se cadastrem** |
| Mandar **lembrete 7 dias antes** do pagamento |
| Mandar **aviso 1 dia depois** se não pagar |
| Marcar o morador como **inadimplente** se não pagar |
| Ter uma área para **reservar** mostrando o que tem disponível |
| Mostrar **horários já reservados** |
| **Impedir** reserva no mesmo lugar, dia e horário |
| **Impedir** morador inadimplente de reservar |
| Só o **responsável** pode reservar |
| Só o **responsável** pode cadastrar visitante |
| Porteiro pode **ver** visitantes cadastrados |
| Porteiro pode **registrar entrada e saída** |
| Síndico pode **postar comunicados** |
| Mostrar **comunicados** no dashboard |
| Síndico pode **aprovar, editar ou cancelar** reservas |
| Morador pode ver **histórico de pagamentos** |
| Morador pode ver **histórico de reservas** |

---

## Requisitos não funcionais

### Segurança

| O que o sistema precisa ter |
|---|
| Dados pessoais guardados com **criptografia** |
| Cada pessoa acessa só o que o seu perfil permite |
| Toda comunicação deve ser **HTTPS** |
| Login seguro para todos os usuários |

### Usabilidade

| O que o sistema precisa ter |
|---|
| Dashboard que funcione no **computador, tablet e celular** |
| Interface **fácil de usar**, sem complicação |
| Reserva feita em **até 3 minutos** sem precisar de ajuda |

### Desempenho

| O que o sistema precisa ter |
|---|
| Dashboard carrega **rápido** |
| Reserva processada em **até 3 segundos** |
| E-mail e SMS saem em **até 5 minutos** |

### Disponibilidade

| O que o sistema precisa ter |
|---|
| Sempre funcionando (exceto manutenção programada) |
| Vários usuários ao mesmo tempo sem travar |

### Conformidade

| O que o sistema precisa ter |
|---|
| Seguir a **LGPD** (proteção de dados pessoais) |

---

## Revisão técnica

### Banco de Dados

O banco (`database/schema.sql`) tem 8 tabelas com ID automático, chaves estrangeiras e regras de verificação.

Pontos bons:
- ID gerado automaticamente com `gen_random_uuid()`
- Status do usuário alinhado às regras (INADIMPLENTE, ADMINISTRADOR)
- Visitante vinculado ao morador que cadastrou
- Índice para evitar reserva duplicada
- Data de criação e atualização em todas as tabelas
- Schema carrega automático com Docker

Melhorias para fazer:
1. Guardar senha com hash (bcrypt), não em texto plano
2. Criar campo para identificar o morador responsável pelo apartamento
3. Criar tabela de notificações para e-mail e SMS

### Documentação

A documentação cobre os problemas e os requisitos. Para melhorar:
- Manter os arquivos `01-Introducao.md` e `02-Requisitos.md` como referência
- Alinhar os diagramas com o que está escrito nos requisitos
- Corrigir erros de digitação nos diagramas

### Modelo geral

O modelo cobre tudo que o sistema precisa: pagamentos, reservas, visitantes, comunicados e controle de acesso. As melhorias listadas devem ser feitas antes de começar a programar o backend.
