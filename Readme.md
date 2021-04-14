# Eltrade CC300 DRIVER API
v0.1.3

Ce Driver permet de communiquer en http avec le module electronique de controle de facturation de la DGI 

### SETUP

Le Driver est un mini serveur http. Il reçoit les commandes par http et les éxécutent sur l'équipement connecté puis, fournis le resultat en response de la requête HTTP.

Spécifications du serveur
- Protocole : HTTP
- Data Format : JSON
- Authentification: Host Headers [NOT YET IMPLEMENTED]
- ENDPOINT : http://localhost:38917

Pour utiliser cette API il faut :
  - Installer le driver 
  - Envoyer des requetes HTTP/REST sur l'adresse de l'API
  - Traiter les reponses 

## USAGE
### L'API fournit trois (03) commandes : 
##### 1 - [CHECK] : Elle permet de vérifier si le périphérique est connecté à l'ordinateur client

    Path: /check
    Method: GET
    Body: N/A
    Response:
      -  Quand l'équipement n'est pas connecté!
            code Http : 503
            corps de la réponse :  
                {"status": "DeviceNotConnected"} 
                
      -  Quand l'équipement  est bien  connecté!
            code Http : 200
            corps de la réponse :  
                {  "status": "Ready" } 
 
##### 2 - [INFO] : Elle permet d'avoir des informations sur l'équipement et sur le contribuable :
NB: Cette requête prends 4s pour s'exécuter 

    Path: /info
    Body: N/A
    Method: GET
    Response:
      -  Quand l'équipement n'est pas connecté!
            code Http : 503
            corps de la réponse :  
                {"status": "DeviceNotConnected"} 
                
      - Quand l'équipement est  connecté!
            code Http : 200
            corps de la réponse :  
                {
                    "NIM": "ED04000623",
                    "IFU": "3201910768821",
                    "TIME": "2020-05-03 13:24:18 +0100 WAT",
                    "COUNTER": "45",
                    "SellBillCounter": "40",
                    "SettlementBillCounter": "0",
                    "TaxA": "0.00",
                    "TaxB": "18.00",
                    "TaxC": "0.00",
                    "TaxD": "18.00",
                    "CompanyName": "BFT",
                    "CompanyLocationAddress": "RUE 12.170 12 IEME ARRONDISSEMENT",
                    "CompanyLocationCity": "COTONOU",
                    "CompanyContactPhone": "61006060",
                    "CompanyContactEmail": "contact@bftgroup.co",
                    "LastConnectionToServer": "2020-05-03 13:23:31 +0100 WAT",
                    "DocumentOnDeviceCount": "45",
                    "UploadedDocumentCount": "45"
                }

##### 3 - [BILL] : Elle permet de créer une facture : 
 Pour créer une facture vous devez respecter le schéma json d'un de l'objet Bill (voir ```bill.spec.json```) 
 
 NB: Utilisez le fichier bill.spec.json avec un validateur comme https://www.jsonschemavalidator.net/ par example pour valider votre json
  
    Path: /bill
    Body: 
        {
            "seller_id": "ABCDEFGHIJKLMN",
            "seller_name": "ABCDEFGHIJK",
            "payments": [{
                    "mode": "V",
                    "amount": 129.5
                }],
            "products": [{
                "label": "ABCD",
                "tax": "B",
                "price": 807.75,
                "bar_code": "ABCDEFG",
                "items": 723.75,
                "specific_tax": 504.5,
                "specific_tax_desc": "ABCDEFGHI",
                "original_price": 80.75,
                "price_change_explanation": "ABCDEFGHIJKLMN"
            }],
            "rt": "FA",
            "rn": "ABCDEFGHIJKLMNOPQRST",
            "buyer_ifu": "ABCDEF",
            "buyer_name": "ABCDEFGHIJKLMNOPQRSTUV",
            "aib": "N/A"
        }
    Method: POST
    Response:
      -  Quand l'équipement n'est pas connecté!
            code Http : 503
            corps de la réponse :  
                {"status": "DeviceNotConnected"} 
     -  Quand l'équipement est pas connecté et que tous s'est bien passé 
            code Http : 200
            corps de la réponse :  
                {"qr_code":"F;ED04000623;NW34NID6ZHANFNMZ2IU7LL3H;3201910768821;20200503135002"}
