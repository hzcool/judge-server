#include <bits/stdc++.h>

using namespace std;

ifstream f1,f2;
void nxtInt(int &x) {
    f1>>x;
    if(f1.fail()) exit(1);
}

//argv1 是用户的输出文件，argv2是输入文件
int main(int argc,char *argv[]) {
    // ofstream f1(argv[1]),f2(argv[2]);
    f1.open(argv[1]);
    f2.open(argv[2]);
    int x,y,z,n;
    f2>>n;
    for(int i=1;i<=n;i++) {
        nxtInt(x);
        nxtInt(y);
        nxtInt(z);
        if(x+y != z || z>100 || z<0) exit(1);
    } 
    char ed;
    f1>>ed;
    exit((int)!f1.fail());
    return 0;
}

//spj代码 数据输入文件 数据输出文件