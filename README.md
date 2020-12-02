## Laboratorio N°2 Sistemas distribuidos Grupo GCR - G
##### Integrantes:

| Nombre  |   Rol|
| ------------ | ------------ |
| Cristian Bernal R.  | 201773026-9   |
|  Raúl Álvarez C. |  201773010-2 |



#### Como correr el código en cada maquina
Dependiendo de la máquina, se deberá ejecutar el comando make correspondiente:

| Comando make | Máquina     | Funcionalidad |
|--------------|-------------|---------------|
| runn         | 10.10.28.81 | Namenode      |
| rund1        | 10.10.28.82 | Datanode1     |
| rund2        | 10.10.28.83 | Datanode2     |
| rund3        | 10.10.28.84 | Datanode3     |
| runc         | cualquiera  | Cliente       |


#### Funcionamiento del Código
Se usó el lenguaje Go junto con RabbitMQ y gRPC para los envios de datos entre los nodos y el cliente

#### Observaciones
Namenode nunca debería morir, si no el sistema muere ya que este trabaja con datos en memoria.  
No hicimos un filtro de duplicado a la hora de subir un mismo libro de nuevo al servidor.  
En informe esta explicado a grandes rasgos como funciona el sistema, pero el detalle a fondo esta en el propio codigo comentado.  
Los archivos a subir en el sistema deben estar previamente almacenados en la carpeta cliente. Inicialmente se dejará subido un solo archivo llamado "ejemplo.pdf".  
Los archivos a descargar serán guardados en la carpeta cliente, en caso de descargar uno con el nombre de un archivo ya existente en la misma, este se sobrescribirá.  

#### Funcionamiento del sistema
Al momento de inicializar el sistema se requiere que al menos un nodo este inicializado con el cliente, dado que este indicará el algoritmo con el cual funcionarán los demás.  
Ya teniendo en un nodo con el algoritmo implementado, el resto de nodos se acoplarán a este replicando la configuración, con lo cual, para implementar otro algoritmo se requiere que esten todos los nodos apagados.  
Cada vez que se inicie una nueva simulación se recomienda borrar los archivos generados en el anterior (chunks y registro) para evitar incongruencias en el sistema o corrupción de archivos.
