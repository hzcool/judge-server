#include <iostream>
using namespace std;
using ll = long long;
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
    cin>>a>>b;
    cout<<qp(a%mod,b)<<endl;
    return 0;
}