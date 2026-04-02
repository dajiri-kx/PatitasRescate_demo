<!DOCTYPE html>
<html lang="es">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>Próximamente - Patitas al Rescate</title>
  <style>
    body {
      margin: 0;
      padding: 0;
      font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
      background-color: #f9f9f9;
      color: #333;
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      height: 100vh;
      text-align: center;
    }

    .container {
      background: white;
      padding: 2rem 3rem;
      border-radius: 12px;
      box-shadow: 0 4px 12px rgba(0,0,0,0.1);
      max-width: 400px;
      width: 90%;
    }

    h1 {
      font-size: 1.8rem;
      margin-bottom: 1rem;
      color:#fe0401;
    }

    p {
      font-size: 1rem;
      margin-bottom: 1.5rem;
    }

    .btn {
      text-decoration: none;
      background-color:#fe0401;
      color: white;
      padding: 0.6rem 1.2rem;
      border-radius: 8px;
      transition: background-color 0.3s ease;
      display: inline-block;
    }

    .btn:hover {
      background-color: #fe0401;
    }

    .icon {
      font-size: 3rem;
      color: #fe0401;
      margin-bottom: 1rem;
    }
  </style>
</head>
<body>
  <div class="container">
    <div class="icon">❌</div>
    <h1>¡Lo sentimos!</h1>
    <p>Ha ocurrido un error. En este momento no se pueden eliminar las mascotas.<br> Por favor, contactenos.</p>
    <a href="contactenos.php" class="btn">← Contactar al personal</a>
  </div>
</body>
</html>
