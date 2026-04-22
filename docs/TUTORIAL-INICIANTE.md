# 🚀 Tutorial Zero-to-Hero: Programando com Inteligência Artificial

Bem-vindo ao **Kortex**! Este guia foi feito especialmente para você que **quer criar programas, sites ou automações incríveis usando IA**, mas talvez nunca tenha programado na vida ou não tenha tempo para estudar códigos complexos.

Se você seguir este passo a passo, em poucos minutos terá um **Analista de Elite** morando no seu computador, pronto para codificar qualquer ideia que você tiver.

---

## 🧠 1. O que é o Kortex?

Imagine que você contratou o melhor programador do mundo (nós o chamamos de Analista de Elite). Em vez de você escrever os códigos, você apenas conversa com ele e dá as ordens.

O **Kortex** é o "escritório" que nós montamos para ele no seu computador. Ele instala as ferramentas de IA, configura tudo sozinho e, o mais importante, **dá memória ao seu agente** e ensina a ele o nosso método de trabalho profissional.

---

## 🛠️ 2. O que você precisa ter instalado?

Antes de começarmos, seu computador precisa ter duas coisas básicas:

1. **Git**: Para baixar o nosso código. ([Baixar Git](https://git-scm.com/downloads))
2. **Go (Golang)**: A linguagem em que o Kortex foi feito (versão 1.24 ou superior). ([Baixar Go](https://go.dev/dl/))

*Dica:* Se você usa Mac, pode instalar tudo abrindo o Terminal e digitando `brew install git go`. No Windows, baixe os instaladores pelos links acima e siga o "Avançar > Avançar > Concluir".

---

## ⚡ 3. Instalando o seu Analista de Elite

Abra o seu **Terminal** (no Mac/Linux) ou o **PowerShell** (no Windows) e digite (ou copie e cole) os comandos abaixo, apertando `Enter` após cada linha:

**Passo 1: Baixar o Kortex**
```bash
git clone https://github.com/fortissolucoescontato-bit/kortex.git
cd kortex
```

**Passo 2: Construir o programa**
```bash
go build -o kortex ./cmd/kortex
```

**Passo 3: Rodar o instalador mágico**
```bash
./kortex install
```

Nesta tela, o Kortex vai analisar o seu computador e te fazer algumas perguntas usando uma interface super bonita. Apenas siga as instruções na tela e escolha qual agente de IA você prefere (por exemplo, o OpenCode ou o Claude Code). O Kortex fará todo o trabalho duro de baixar e configurar!

---

## 🎮 4. Como usar na prática? (O Jeito Básico)

Com tudo instalado, você não precisa mais se preocupar com configurações. Para começar a trabalhar, basta abrir o seu Terminal na pasta onde você quer criar o seu projeto e chamar o seu agente!

Por exemplo, se você instalou o **OpenCode**, digite:

```bash
opencode
```

O agente vai iniciar uma conversa com você. Tente pedir algo simples como:
> *"Crie um script em Python para organizar minhas fotos por data."*

Ele vai escrever o código para você. Mas e se o projeto for grande, como um site completo ou um sistema de vendas? É aí que entra o verdadeiro poder do Kortex.

---

## 🏗️ 5. O Superpoder Profissional: O Método SDD

A maioria das pessoas usa a IA de forma errada: pede um sistema inteiro em uma única frase. O resultado? Código bagunçado que quebra no dia seguinte.

O Kortex usa uma metodologia chamada **SDD (Spec-Driven Development)**, ou Desenvolvimento Orientado a Especificações. Em vez de programar cegamente, o seu Analista de Elite vai primeiro **explorar, propor, desenhar e só depois programar**, exatamente como uma equipe de engenharia da Nexo-Fortis faria.

E o melhor: você controla tudo isso através de comandos mágicos.

### Os Comandos "Piloto Automático" (Recomendado para Iniciantes)

Se você quer que a IA faça o trabalho pesado de forma estruturada sem que você precise gerenciar cada etapa, use estes dois meta-comandos:

1. **`/sdd-new <o que você quer>`**: Inicia um projeto do zero.
   *Exemplo:* `/sdd-new Crie um sistema de login para minha padaria`
   *O que ele faz:* O agente vai investigar, criar uma proposta de arquitetura, escrever as especificações, quebrar tudo em tarefas menores e te pedir autorização para começar a programar.

2. **`/sdd-continue`**: O comando que você mais vai usar. 
   Sempre que o agente terminar uma fase (por exemplo, terminou de escrever a especificação), basta digitar `/sdd-continue` para ele passar para a próxima fase (como desenhar a arquitetura ou programar o código).

### Os Comandos "Engenheiro Chefe" (O Fluxo Completo)

Se você quer ter controle cirúrgico sobre cada etapa do processo, você pode chamar cada fase individualmente. Este é o ciclo de vida completo de um software no Kortex:

| Comando | Para que serve? | O que o agente faz? |
|---------|-----------------|---------------------|
| **`/sdd-init`** | **O Primeiro Passo** | Prepara o seu projeto. Descobre quais linguagens você está usando e liga o modo de testes rígidos (TDD) se necessário. |
| **`/sdd-explore`** | **A Investigação** | Você tem uma ideia e quer saber se é possível? Use este comando. O agente vai ler o código e comparar soluções (Ex: `/sdd-explore como podemos colocar um chat aqui?`). |
| **`/sdd-propose`** | **A Proposta** | Cria um documento com o escopo, a intenção e o plano de rollback (como reverter caso dê erro). |
| **`/sdd-spec`** | **As Regras** | Escreve os requisitos técnicos e as histórias de usuário (o que o sistema *deve* fazer). |
| **`/sdd-design`** | **A Arquitetura** | Define quais bancos de dados e padrões de código serão usados na solução. |
| **`/sdd-tasks`** | **O Checklist** | Quebra todo o planejamento em micro-tarefas (ex: 1. criar botão, 2. configurar banco). |
| **`/sdd-apply`** | **A Mão na Massa** | O agente finalmente começa a escrever o código real, seguindo estritamente as tarefas criadas no passo anterior. |
| **`/sdd-verify`** | **O Controle de Qualidade** | Verifica se o código escrito realmente funciona e atende às especificações. |
| **`/sdd-archive`** | **O Fechamento** | Marca a funcionalidade como concluída, sincroniza a documentação e limpa o ambiente. |

**Dica de Ouro:** Não quer digitar tudo isso passo a passo? Use o **`/sdd-ff <nome-do-projeto>`** (Fast-Forward). Ele roda a proposta, a especificação, o design e as tarefas tudo de uma vez na sequência, parando apenas antes de começar a codificar!

---

## 💾 6. A Memória Implacável: O Engram

A maior frustração com IAs normais (como o ChatGPT) é que se você fechar a janela, a IA esquece tudo o que vocês fizeram.

No ecossistema Kortex, o seu agente é integrado ao **Kortex Engram** (uma memória persistente de longo prazo). Todo comando `/sdd` que você executa salva automaticamente decisões de arquitetura, padrões e correções de bugs.

Quando você voltar a trabalhar no projeto no dia seguinte, ou até daqui a um ano, o agente vai silenciosamente buscar no banco de dados e vai **lembrar de todas as decisões que vocês tomaram no passado!** 

*(Nota: Graças ao protocolo de Auditoria de Contexto, o Kortex sempre lê os resumos das últimas sessões antes de iniciar qualquer fase do SDD, garantindo que ele nunca "quebre" ou desfaça algo que já estava funcionando).*

---

## 🎉 Conclusão

Pronto! Você agora tem um ambiente de engenharia profissional rodando na sua máquina. A partir de agora, **sua única limitação é a sua criatividade**.

**Resumo para o sucesso:**
1. Comece ideias grandes com `/sdd-new`.
2. Avance as etapas do planejamento com `/sdd-continue`.
3. Deixe o **Engram** cuidar da memória do projeto.

Dúvidas ou problemas técnicos? Envie uma issue (mensagem) lá no nosso GitHub e estaremos prontos para ajudar. Mãos à obra!
