<!-- factura.php -->
<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="utf-8">
    <title>Factura - Patitas al Rescate</title>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css" rel="stylesheet">
    <style type="text/css">
        body { margin-top: 20px; background-color: rgba(235, 238, 255, 1) !important; }
        .invoice-aside { background-color: #ededed; display: flex; flex-direction: column; padding: 50px 33px; min-width: 300px; }
        .invoice-content { background-color: #fff; padding: 50px 33px; flex: 1 1 0%; }
        .invoice-title { font-size: 2.5rem; font-weight: 300; }
        .invoice-order { text-align: right; }
        .invoice-details, .invoice-summary { width: 100%; font-size: 1.1rem; margin-bottom: 50px; }
        .invoice-summary { border-top: 1px solid #d9d9d9; font-size: 1.3rem; }
        .invoice-summary th { padding-top: 20px; font-weight: 600; width: 20%; }
        .invoice-summary th.total { width: 60%; font-size: 1.8rem; text-align: right; }
        .invoice-summary .total-value { text-align: right; font-size: 2.5rem; }
        .invoice-payment-details { border: 1px solid #d9d9d9; padding: 20px; }
    </style>
</head>
<body>
    <?php include 'header.php'; ?>

    <div class="container">
        <div class="row invoice">
            <div class="col-md-3 invoice-aside">
                <div class="invoice-person">
                    <span class="name">María Rodríguez</span>
                    <span class="position">Veterinaria</span>
                    <span>San José, Costa Rica</span>
                </div>
                <div class="invoice-person mt-5">
                    <span class="name">Juan Corredor</span>
                    <span class="position">Cliente</span>
                    <span>Heredia, Costa Rica</span>
                </div>
            </div>
            <div class="col-md-9 invoice-content">
                <div class="row invoice-header">
                    <div class="col-6 invoice-title">
                        <span>Factura</span>
                    </div>
                    <div class="col-6 invoice-order">
                        <span class="invoice-number">Número: FAC-2025-0001</span>
                        <span class="invoice-date">Fecha: 17 de marzo de 2025</span>
                    </div>
                </div>
                <div class="row">
                    <div class="col-md-12">
                        <table class="invoice-details">
                            <thead>
                                <tr>
                                    <th style="width:60%">Descripción</th>
                                    <th class="amount" style="width:40%">Monto</th>
                                </tr>
                            </thead>
                            <tbody>
                                <tr>
                                    <td>Consulta veterinaria</td>
                                    <td class="amount">₡25,000</td>
                                </tr>
                                <tr>
                                    <td>Baño y corte de pelo</td>
                                    <td class="amount">₡15,000</td>
                                </tr>
                                <tr>
                                    <td>Vacunación</td>
                                    <td class="amount">₡10,000</td>
                                </tr>
                            </tbody>
                        </table>
                        <table class="invoice-summary">
                            <thead>
                                <tr>
                                    <th>Subtotal</th>
                                    <th>IVA (13%)</th>
                                    <th class="total">Total</th>
                                </tr>
                            </thead>
                            <tbody>
                                <tr>
                                    <td class="amount">₡50,000</td>
                                    <td class="amount">₡6,500</td>
                                    <td class="amount total-value">₡56,500</td>
                                </tr>
                            </tbody>
                        </table>
                    </div>
                </div>
                <div class="row">
                    <div class="col-md-12">
                        <div class="invoice-payment-details">
                            <p><b>Método de pago:</b> Tarjeta</p>
                            <p><b>Tipo de tarjeta:</b> Visa</p>
                            <p><b>Número de verificación:</b> 1234567890</p>
                        </div>
                    </div>
                </div>
                <div class="row invoice-footer mt-4">
                    <div class="col-md-12 text-end">
                        <button class="btn btn-primary">Pagar ahora</button>
                    </div>
                </div>
            </div>
        </div>
    </div>

    <?php include 'footer.php'; ?>

    <script src="https://code.jquery.com/jquery-1.10.2.min.js"></script>
    <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/js/bootstrap.bundle.min.js"></script>
</body>
</html>
