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
      color:#277e1c;
    }

    p {
      font-size: 1rem;
      margin-bottom: 1.5rem;
    }

    .btn {
      text-decoration: none;
      background-color: #277e1c;
      color: white;
      padding: 0.6rem 1.2rem;
      border-radius: 8px;
      transition: background-color 0.3s ease;
      display: inline-block;
    }

    .btn:hover {
      background-color: #277e1c;
    }

    .icon {
      font-size: 3rem;
      color: #277e1c;
      margin-bottom: 1rem;
    }
  </style>
</head>
<body>
  <div class="container">
    <div class="icon">✉️</div>
    <h1>¡Mensaje Enviado Exitosamente!</h1>
    <p>Muchas gracias por tu mensaje.<br> ¡Pronto nos pondremos en contacto con usted!</p>
    <a href="home.php" class="btn">← Volver al Inicio</a>
  </div>
</body>
</html>
