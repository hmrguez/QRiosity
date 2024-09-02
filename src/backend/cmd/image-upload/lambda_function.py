import json
import boto3
import base64
import os
import mimetypes

s3 = boto3.client('s3')

def lambda_handler(event, context):
    try:
        # Extract the body from the event, which contains the base64-encoded image data
        body = event['body']

        print("Event: ", event)

        # If the body is base64 encoded, decode it
        if event.get('isBase64Encoded', False):
            image_data = base64.b64decode(body)
#         else:
#             return {
#                 'statusCode': 400,
#                 'body': json.dumps({'error': 'File data must be base64 encoded.'})
#             }

        print("Event: ", event)
        print("Body: ", body)

        # Get the Content-Type header from the event
        content_type = event['body']['headers']['Content-Type']

        print("Content-Type: ", content_type)



        # Determine the content type of the image from headers
#         content_type = event['body'].get('Content-Type')
        
        # Validate the content type
        if content_type not in ['image/jpeg', 'image/jpg', 'image/png', 'image/webp']:
            return {
                'statusCode': 400,
                'body': json.dumps({'error': 'Unsupported file type. Only jpeg, jpg, png, and webp are allowed.'})
            }
        
        # Generate a unique key for the image
        file_extension = mimetypes.guess_extension(content_type)
        image_key = f"uploads/{context.aws_request_id}{file_extension}"
        
        # Upload the image to S3 with correct content type
        s3.put_object(Bucket=os.environ['BUCKET_NAME'], 
                      Key=image_key, 
                      Body=image_data, 
                      ContentType=content_type)

        return {
            'statusCode': 200,
            'body': json.dumps({'message': 'Image uploaded successfully', 'image_key': image_key})
        }
    except Exception as e:
        return {
            'statusCode': 500,
            'body': json.dumps({'error': str(e)})
        }