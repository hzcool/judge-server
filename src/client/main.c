#include <stdio.h>

typedef long long ll;
const int mod = 1e9+7;
ll qp(ll a,ll b) {
    ll ans = 1;
    while(b) {
        if(b&1) ans = ans * a % mod;
        b>>=1;
        a = a*a%mod;
    }
    return ans;
}
int main() {
    ll a,b;
    scanf("%lld%lld",&a,&b);
    printf("%lld\n",qp(a%mod,b));
    return 0;
}