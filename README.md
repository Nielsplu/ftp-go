# Projet Ftp 

FILMONT Félix<br>
AIGNELOT Youenn<br>
PLU Niels<br>

### Mise en place

Installation les packets :

```bash
go mod download
```

Lancer le serveur :

```bash
go build -o server cmd/server/main.go
./server
```

Lancer le/les client :

```bash
go build -o client cmd/client/main.go
./client
```

### Fonctionnalités : 

- Commande `List`
- Commande `Cd`
- Commande `Get`
  - Cette commande se décline en deux implémentation<br>
    Une pour les petits fichiers qui les télécharges d'un coup<br>
    Et une pour les gros fichier qui sont envoyés chunk par chunks
- Commande `End`
- Commande `Hide`
- Commande `Reveal`
- Commande `Terminate`

- Timout d'une minute pour les connection inactive.
- Servir une arborescence
- Sécurité pour ne pas échapper le fichier data __server__

### Race conditions

Notre infrastructure __(server)__ se base sur des **Stoppers**, un design patern 
implémenté par nos soins. Un stopper a une liste de **stopper enfant**, il peut 
leur envoyer un message pour qu'ils se termine gracieusement, et attendre la fin 
de tous ses enfants. 

Les stoppers sont utilisé par exemple dans les `Accept Loop`, un enfant par port ouvert.<br>
Également par les `Reader` de chaque connection 
__(chaque connection est divisé en un reader et un writer)__ le reader est le maître, il possède un 
stopper pour stopper le writer et attendre sa fin.


## Essayer dans le navigateur

Le vrai serveur et le vrai client tournent en WebAssembly sur
[le portfolio](https://nielsplu.github.io/portfolio/) (carte « Serveur et
client FTP en Go », bouton « Tester ici ») : ils sont compilés tels quels et
reliés par le réseau loopback en mémoire de Go, seule la vue Bubbletea est
remplacée par la page.

Pour compiler le binaire wasm :

```
cd cmd/wasm
GOOS=js GOARCH=wasm go build -o ftp.wasm .
```
