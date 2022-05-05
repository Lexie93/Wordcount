Istruzioni per avviare l'applicazione

Dopo aver generato gli eseguibili (go build worker.go, go build master.go)

Seguire uno dei seguenti due metodi:

A) lanciare lo script "launcher.sh" indicando come argomenti nel seguente ordine:
1) il numero di worker che si desidera utilizzare
2) un elenco di file e/o cartelle (verranno automaticamente aggiunti all'elenco dei file quelli nelle cartelle indicate)
Es.: $ ./launcher.sh 4 file1.txt file2.txt cartella file3.txt

B) lanciare singolarmente ogni worker ed il master con i seguenti argomenti
- per i worker:
1) porta da utilizzare (in modo sequenziale 1234 -> 1235 -> ...)
Es.: $ ./worker 1234
- per il master (dopo aver avviato i worker):
1) il numero di worker che si desidera utilizzare
2) un elenco di file e/o cartelle (verranno automaticamente aggiunti all'elenco dei file quelli nelle cartelle indicate)
Es.: $ ./master 4 file1.txt file2.txt cartella file3.txt
