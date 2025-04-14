import argparse
import tiktoken
import json
from google import genai
from google.genai.types import HttpOptions

import sys

parser = argparse.ArgumentParser(description="Tokenize text using Google GenAI API")

def process_gemini(texts: list[str]) -> list[int]:
    client = genai.Client(http_options=HttpOptions(api_version="v1"))
    response = client.models.compute_tokens(
        model="gemini-2.0-flash-001",
        contents=texts
    )

    res = []
    for tokens_info in response.tokens_info:
        res.append(len(tokens_info.token_ids))
    
    return res

def process_openai(texts: list[str]) -> list[int]:
    enc = tiktoken.encoding_for_model("gpt-3.5-turbo")
    res = []
    for text in texts:
        tokens = enc.encode(text)
        res.append(len(tokens))
    return res

def main():
    parser.add_argument("-p", "--platform", type=str, required=True, help="Platform name: openai or gemini")
    parser.add_argument("-t", "--texts", type=str, nargs="+",required=True, help="Texts to tokenize, separated by space")
    args = parser.parse_args()

    if args.platform == "gemini":
        try:
            token_counts = process_gemini(args.texts)
            print(json.dumps({"data": token_counts}))
        except Exception as e:
            print(f"Error processing Gemini: {e}", file=sys.stderr)
            sys.exit(1)
    
    elif args.platform == "openai":
        try:
            token_counts = process_openai(args.texts)
            print(json.dumps({"data": token_counts}))
        except Exception as e:
            print(f"Error processing OpenAI: {e}", file=sys.stderr)
            sys.exit(1)
    else:
        print("Unsupported platform. Please use 'openai' or 'gemini'.", file=sys.stderr)
        sys.exit(1)
    


if __name__ == "__main__":
    main()
