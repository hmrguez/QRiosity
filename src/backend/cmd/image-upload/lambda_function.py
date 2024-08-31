import json
import boto3
import base64
import random
import string

s3 = boto3.client('s3')

def is_image(content_type):
    return content_type in ['image/jpeg', 'image/jpg', 'image/png', 'image/webp']

def lambda_handler(event, context):
    try:
        # Get the file from the event
        file_content = base64.b64decode(event['body'])
        content_type = event['headers']['Content-Type']
        file_size = len(file_content)

        # Check if the file is an image
        if not is_image(content_type):
            return {
                'statusCode': 400,
                'headers': {
                    'Access-Control-Allow-Origin': '*',  # Allow CORS
                    'Access-Control-Allow-Methods': 'POST, OPTIONS',
                    'Access-Control-Allow-Headers': 'Content-Type'
                },
                'body': json.dumps('File type not supported. Please upload a JPEG, JPG, PNG, or WEBP image.')
            }

        # Check file size (less than 5MB)
        if file_size > 5 * 1024 * 1024:
            return {
                'statusCode': 400,
                'headers': {
                    'Access-Control-Allow-Origin': '*',  # Allow CORS
                    'Access-Control-Allow-Methods': 'POST, OPTIONS',
                    'Access-Control-Allow-Headers': 'Content-Type'
                },
                'body': json.dumps('File size exceeds 5MB limit.')
            }

        # Generate a random 15-letter string for the filename
        file_extension = content_type.split('/')[-1]
        random_filename = ''.join(random.choices(string.ascii_letters, k=15)) + '.' + file_extension

        # Define the S3 bucket name and file key
        bucket_name = 'your-s3-bucket-name'
        file_key = random_filename

        # Upload the file to S3
        s3.put_object(Bucket=bucket_name, Key=file_key, Body=file_content, ContentType=content_type)

        # Generate the URL of the uploaded file
        file_url = f'https://{bucket_name}.s3.amazonaws.com/{file_key}'

        # Return the URL as the response body
        return {
            'statusCode': 200,
            'headers': {
                'Access-Control-Allow-Origin': '*',  # Allow CORS
                'Access-Control-Allow-Methods': 'POST, OPTIONS',
                'Access-Control-Allow-Headers': 'Content-Type'
            },
            'body': json.dumps({'file_url': file_url})
        }

    except Exception as e:
        return {
            'statusCode': 500,
            'headers': {
                'Access-Control-Allow-Origin': '*',  # Allow CORS
                'Access-Control-Allow-Methods': 'POST, OPTIONS',
                'Access-Control-Allow-Headers': 'Content-Type'
            },
            'body': json.dumps(f'Error processing file: {str(e)}')
        }