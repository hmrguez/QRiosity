import boto3
import base64
import uuid
import mimetypes
import json

s3 = boto3.client('s3')

def lambda_handler(event, context):
    # Get the base64-encoded image and MIME type from the event
    print(event)

#     print("Event: ", event)

    event = json.loads(event['body'])

#     print("Body: ", event)

    image_base64 = event['image']
    mime_type = event['mimeType']

#     print("Image: ", image_base64)


    # Decode the image
    image_data = base64.b64decode(image_base64)

    # Determine the file extension based on MIME type
    extension = mimetypes.guess_extension(mime_type)
    if extension is None:
        raise ValueError(f"Unsupported MIME type: {mime_type}")


    # Generate a unique filename with the correct extension
    filename = f"{uuid.uuid4()}{extension}"


    # Specify your S3 bucket name
    bucket_name = 'roadmap-images'

    try:
        # Upload the image to S3
        s3.put_object(Bucket=bucket_name, Key=filename, Body=image_data, ContentType=mime_type)

        # Generate the S3 URL
        s3_url = f"https://{bucket_name}.s3.amazonaws.com/{filename}"

        body = json.dumps({
            'message': f'Successfully uploaded {filename} to {bucket_name}',
            'url': s3_url
        })

        return {
            'statusCode': 200,
            'body': body
        }
    except Exception as e:

        body = json.dumps({
            'message': f'Error uploading image: {str(e)}',
            'url': None
        })

        return {
            'statusCode': 500,
            'body': body
        }