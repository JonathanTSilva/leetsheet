import json

# Paste your code between the triple quotes
code = """class Solution:
    def longestDecomposition(self, S: str) -> int:
        # Initialize result counter and left/right substring accumulators
        res, l, r = 0, '', ''
        
        # Iterate through the string from both ends simultaneously
        # i iterates from left to right, j iterates from right to left
        for i, j in zip(S, S[::-1]):
            # Build left substring by appending current character from left
            # Build right substring by prepending current character from right
            l, r = l + i, j + r
            
            # When left and right substrings match, we found a valid chunk pair
            if l == r:
                # Increment result count and reset substring accumulators
                res, l, r = res + 1, '', ''
        
        return res
"""

# Convert to a JSON-safe string
escaped = json.dumps(code)

# Print it
print(escaped)

# Optionally copy to clipboard (so you can just paste it into your JSON file)
try:
    pyperclip.copy(escaped)
    print('\n✅ Copied to clipboard!')
except Exception:
    print('\n⚠️ pyperclip not installed or clipboard unavailable')
