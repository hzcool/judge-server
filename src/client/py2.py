mod = int(1e9+7)
def qp(a,b):
    ans = 1
    while(b):
        if (b&1):
            ans = ans * a % mod
        a = a * a % mod
        b >>= 1
    return ans

str = raw_input()
arr = str.split(" ")
a = int(arr[0])
b = int(arr[1])
print(qp(a,b))
