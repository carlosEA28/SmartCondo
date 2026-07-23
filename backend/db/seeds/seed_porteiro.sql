INSERT INTO Usuario (nome, email, senha, telefone, tipo)
VALUES
  ('Porteiro Plantão', 'porteiro@condominio.com', '$2a$10$tfRFTF9BvuwuyNHzcvKbwebfK4Ay3fbw4atrnjR7mhZ7dX.iTsnNy', '11999999999', 'PORTEIRO'),
  ('Síndico Geral', 'sindico@condominio.com', '$2a$10$tfRFTF9BvuwuyNHzcvKbwebfK4Ay3fbw4atrnjR7mhZ7dX.iTsnNy', '11977777777', 'SINDICO')
ON CONFLICT (email) DO NOTHING;
