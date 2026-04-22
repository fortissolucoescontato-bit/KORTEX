# Antigravity IDE: Workaround para SDD com Agente Único

## O Problema
O Antigravity IDE atualmente não suporta a invocação nativa de subagentes em segundo plano (multi-threading de LLMs). Todas as fases do Desenvolvimento Orientado a Especificações (SDD) devem ser executadas sequencialmente na mesma conversa.
Isso gera um alto risco de **degradação de contexto e alucinações**, pois o LLM começa a misturar instruções de skills anteriores e perde o fio da arquitetura.

## A Solução: Máquina de Estados Baseada em Artefatos
Para implementar uma integração precoce sem depender de APIs externas e **sem afetar a arquitetura multi-agente original de projetos como Cursor ou OpenCode**, aplicamos um padrão de Máquina de Estados apoiada estritamente no File System local.

### Regras a Injetar (Específicas para o Antigravity)

1.  **Troca Estrita de Papéis (Role Switching)**
    O orquestrador deve anunciar a mudança de fase e carregar a skill correspondente (`SKILL.md`) em seu contexto, ignorando temporariamente as diretrizes anteriores.

2.  **File-System como Memória (Save State)**
    Ao terminar uma fase (ex: `sdd-propose`), o Antigravity está PROIBIDO de avançar sem antes salvar a saída completa em um arquivo físico (ex: `.sdd/proposta.md`). O chat NÃO é um meio de armazenamento confiável.

3.  **Amnésia Controlada (Load State)**
    Ao iniciar a fase seguinte (ex: `sdd-spec`), o Antigravity NÃO DEVE confiar em seu histórico de chat. Sua primeira ação obrigatória é usar a ferramenta de leitura (`Read`) para carregar o arquivo gerado no passo anterior. Isso refresca o contexto exato necessário para a fase atual.

4.  **Uso Correto do Engram**
    O Engram (`mem_save`) é preservado UNICAMENTE para registrar decisões arquitetônicas globais, convenções e correções de bugs. NÃO deve ser usado para salvar o estado intermediário de um SDD em andamento (para isso servem os arquivos `.sdd/*.md`).

## Conclusão
Este workaround permite ter o SDD funcional no Antigravity hoje mesmo, operando sob um modo de "Simulação de Thread Única", mantendo intactas a limpeza e a modularidade do repositório original do Kortex.
