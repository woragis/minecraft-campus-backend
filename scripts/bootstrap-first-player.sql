-- Cria o primeiro jogador (fundador) que pode enviar convites.
-- Substitua o UUID e username pelos valores reais do Minecraft.

INSERT INTO players (
    id,
    minecraft_uuid,
    username,
    status,
    trust_score,
    sponsor_score,
    created_at,
    updated_at
) VALUES (
    '11111111-1111-1111-1111-111111111111',  -- troque pelo UUID interno desejado
    '00000000-0000-0000-0000-000000000001',  -- troque pelo UUID Minecraft real
    'Fundador',
    'active',
    100,
    100,
    now(),
    now()
) ON CONFLICT (minecraft_uuid) DO NOTHING;
