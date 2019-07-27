package handlers

import (
	"fmt"
	"github.com/google/tink/go/keyset"
	"github.com/google/tink/go/signature"
	"log"
)

func Verification(sinParams string,sinByte []byte) bool{

	// 初始化 签名验证器, 我们使用的是 SHA256WithECDSA
	khPriv, err := keyset.NewHandle(signature.ECDSAP256KeyTemplate())
	if err != nil {
		log.Fatal(err)
		return false
	}

	s, err := signature.NewSigner(khPriv)
	if err != nil {
		log.Fatal(err)
		return false
	}

	a, err := s.Sign([]byte(sinParams))
	if err != nil {
		log.Fatal(err)
		return false
	}
	// 获取公钥，   TMD 获取公钥地址，java 和 golang 的api 完全不同
	// admob 公钥服务器 获取地址在代码里面写死了
	// 就是这个"type.googleapis.com/google.crypto.tink.EcdsaPrivateKey"
	// 并且 golang 没有 测试公钥列表
	// TODO 此处可能存在问题， 因为ssv回调参数 key_id 并没有使用，而JAVA办admob返回公钥列表是有多个，暂时还不清楚如何获取的指定公钥
	// TODO 也有可能由于不同语言调用api不同，返回的公钥个数不一样
	khPub, err := khPriv.Public()
	if err != nil {
		log.Fatal(err)
		return false
	}

	// 使用公钥初始化验证器
	v, err := signature.NewVerifier(khPub)

	// 拿公钥加密后的签名和url上的签名进行比对
	if err := v.Verify(a, sinByte); err != nil {
		log.Fatal("signature verification failed")
		return false
	}
	// 验证通过
	fmt.Println("Signature verification succeeded.")
	return true
}

/*
签名方法
*/
func Sin(sinParams string)(signatureStr string, err error){
	// 初始化 签名验证器, 我们使用的是 SHA256WithECDSA
	khPriv, err := keyset.NewHandle(signature.ECDSAP256KeyTemplate())
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	s, err := signature.NewSigner(khPriv)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	a, err := s.Sign([]byte(sinParams))

	return string(a), nil
}
