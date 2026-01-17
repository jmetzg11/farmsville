import json
import socket
from django.http import JsonResponse
from django.views.decorators.csrf import csrf_exempt
from django.views.decorators.http import require_POST


PRINTER_IP = "192.168.2.200"
PRINTER_PORT = 9100
TEMPLATE_SLOT = "001"


@csrf_exempt
@require_POST
def print_label(request):
    try:
        data = json.loads(request.body)
        product = data.get('product', '')
        customer = data.get('customer', '')

        # Build P-touch Template command
        # ^II = Initialize
        # ^TS001 = Select template slot 001
        # ^ONobj1\x00...\x00 = Object 1 data
        # ^ONobj2\x00...\x00 = Object 2 data
        # ^FF = Form feed (print and cut)
        command = f"^II^TS{TEMPLATE_SLOT}^ONobj1\x00{product}\x00^ONobj2\x00{customer}\x00^FF"

        # Send to printer via socket
        with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
            s.settimeout(5)
            s.connect((PRINTER_IP, PRINTER_PORT))
            s.sendall(command.encode('utf-8'))

        return JsonResponse({'success': True, 'message': f'Printed: {product} - {customer}'})

    except socket.timeout:
        return JsonResponse({'success': False, 'error': 'Printer connection timed out'}, status=500)
    except socket.error as e:
        return JsonResponse({'success': False, 'error': f'Printer connection failed: {str(e)}'}, status=500)
    except Exception as e:
        return JsonResponse({'success': False, 'error': str(e)}, status=500)


@csrf_exempt
@require_POST
def print_labels_batch(request):
    """Print multiple labels at once"""
    try:
        data = json.loads(request.body)
        labels = data.get('labels', [])

        if not labels:
            return JsonResponse({'success': False, 'error': 'No labels provided'}, status=400)

        # Build commands for all labels
        commands = []
        for label in labels:
            product = label.get('product', '')
            customer = label.get('customer', '')
            command = f"^II^TS{TEMPLATE_SLOT}^ONobj1\x00{product}\x00^ONobj2\x00{customer}\x00^FF"
            commands.append(command)

        # Send all commands to printer
        with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
            s.settimeout(10)
            s.connect((PRINTER_IP, PRINTER_PORT))
            for command in commands:
                s.sendall(command.encode('utf-8'))

        return JsonResponse({'success': True, 'message': f'Printed {len(labels)} labels'})

    except socket.timeout:
        return JsonResponse({'success': False, 'error': 'Printer connection timed out'}, status=500)
    except socket.error as e:
        return JsonResponse({'success': False, 'error': f'Printer connection failed: {str(e)}'}, status=500)
    except Exception as e:
        return JsonResponse({'success': False, 'error': str(e)}, status=500)
