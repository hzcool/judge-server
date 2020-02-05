import java.util.*;
public class Main{
    public static void main(String[] args){
        Scanner s = new Scanner(System.in);
        long mod = 1000000007;
        long a = s.nextLong() % mod;
        long b = s.nextLong();
        long ans = 1;
        while(b>0) {
            if((b&1)>0) ans = ans*a%mod;
            b /= 2;
            a = a*a%mod;
        }
         System.out.println(ans);
    }
}