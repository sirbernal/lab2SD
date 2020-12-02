## Laboratorio N°2 Sistemas distribuidos Grupo GCR - G
##### Integrantes:

| Nombre  |   Rol|
| ------------ | ------------ |
| Cristian Bernal R.  | 201773026-9   |
|  Raúl Álvarez C. |  201773010-2 |



#### Como correr el código en cada maquina
Dependiendo de la máquina, se deberá ejecutar el comando make correspondiente:

|  Comando make | Máquina   |
| ------------ | ------------ |
|   | 10.10.28.81 |
|    | 10.10.28.82  |
|  |   10.10.28.83 |
|   | 10.10.28.84  |


#### Funcionamiento del Código
Se usó el lenguaje Go junto con RabbitMQ y gRPC para los envios de datos entre los nodos y el cliente

#### Observaciones
Namenode nunca debería morir, si no el sistema muere ya que este trabaja con datos en memoria.  
No hicimos un filtro de duplicado a la hora de subir un mismo libro de nuevo al servidor.
En informe esta explicado a grandes rasgos como funciona el sistema, pero el detalle a fondo esta en el propio codigo comentado
