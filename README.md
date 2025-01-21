# SAE Emulateur

## Description

L'objectif de ce projet est de réaliser un émulateur de SAE, étais de réaliser un emulateur risck 5 en go.

Mais je suis rester bloqué a l'instruction `SH: Store halfword` et je n'ai pas pu continuer, donc il ne fonctionne pas.

C'est domage car j'ai compris le fonctionnement de risck 5 et j'ai réussi a faire fonctionner les autres
instructions et le décodage des instructions.

Le mode pas a pas + commande de debug fonctionne correncement et temps que le programme ne fait pas de `SH` il
fonctionne.

Mes fichier de test passe touts sauf instruction_test.go

## Executé un Livrable

```bash
docker build -f Dockerfile_<LIVRABLE> -t sae-emulateur .
docker run --rm -v $(pwd):/data sae-emulateur <ARGS>
```

## Exemple

```bash
docker build -f Dockerfile_1 -t sae-emulateur .
docker run --rm -v $(pwd):/data sae-emulateur -h
```