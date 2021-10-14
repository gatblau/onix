/*
  Onix Config Manager - Pilot
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package core

import (
	"crypto/aes"
	"crypto/cipher"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ProtonMail/gopenpgp/v2/helper"
	"os"
	"time"
)

var (
	k = "819da5fa6489428f9b95780d5f5d740d651b50e21c99b33101eceeb37a5c8850"
	i = "f9e02e9fce1a498a34b08e67"
	a = "233937abe973e03080d18ab92bec0c37ceca093e9301902fd4e7a1f2b910e106dccc9be77b54cd413ec1174c348705bc49b43debc7e0897f2edddaab79faf57219809853c1d864c3c56d31a33d1aa9fb8e6f3b13cdf70cbd63d1c80d915a0747483bd75424b6fd2d27201ce17589ba85186d681f08fabf3b174c681c310681c12b5cb9007d18f6a7ce221691f15b0e11acb5eaef90f9dfa0152e02e0f08266beb123f5a37949dbb6f42ac55b0a4178fa2208e4106a8844185349c19040eb88783f927339c25d2729dde9473fbd4dadc01159624e374f1fd19cede820b8eb3861e22885dd3cf1d423a95badf17c08dfda1c79c01c573a65e03c159c2a818048444cc719affd3615ec18cfc6df91b9e3a7c22cb8af136178e4096da9d20bb5769c45d64813aad7cdb0ba94a302dd2cd7cd68013fe27ba372638e35d4bd2439e15c4e56621e1f0174577b764386ba78eeff07e749431a88e2100f8ba1f7ab164d1866387fbb4726a7820a0d290af14ecd7ce7ee819a874e19c1e77328109d91bfc2dfc508cad3d7262592925ad2bb8a3f9d5032fb21c762e78d9a5d90fe0ec56c8358cc1b5d1707465531ecc19e35428cd43c68e3134660e9c80f062b6d7cf49b0b57ae550baf0e924a4a805873e549a81260e393606cfba443882cb0abd5fc0935de3a878a9bc282294cf7d61d708e68b3280726ffad75e4923f1ef2130988d4b148e26b84187fd90d573f46aee8e6c00e247e42010ba7a005c23c7a837534d6ea38751f48858b5d4a5f27c14d704913b87f3fd72103ddf953171bc2b113d1726a8ba7e1304070931bd6159dcab314f2a9d6352c2062043c3a23e167ebd0d494f23aa9968adfbf4b8f6b67ab971ec555e66414a864d0fdb2e255480e79ee2f4650bf2d916b09ba9d19c1820d870a991b161688051c304db6a3976007aec764feaf1e31c350e11e7e7bd862dfc406b4bfeb6280a082ee49308f1593f5b1002d3920066efb7d3a0c1fcffc07f4535db754ab28ef3381ed83eed01772632897e4ed3cb73b1c77bf17b0be5fa7817a8f2201592f8f433a7eb39e91ecdd24d596c1ba50ead4e9107297e93f54e12dc998c80d56ed1128a8b331172a5e28cf52273101d289d0519e2b09d0755699295150be8e87d00794803b594342ce2fea87e866b4cfd0b609285e8b1d3778fd0f44a534c5b6a56d2a4bcc9203fcd5d4a1d8a0d5775b2e705b39462c39bd75151c6520eac8fe3508d8543f7abed7e738a4f1e17a73e3ef3bb1fbb4027605b48a994ad2e8c61428ce668c767996dfb07233c518bcdca68b25e3e6151d955efe3a21cdc56468b9db87f33e0da22f739e083a2f02ecf97dd879a33c623fdaacf10892f1955273ebe1e9d7df927e7a64777c60bed05e0395fdd0cf02fe93c408113b3483ecc09cad381b88c8aefbd831420bf710f3308649b95fe977013491f117d6281fbe703def86fba362984328adf9a5d22da3f18d7e5e968e8b8afd5da12a7516915f7f74d4ff7c7ff1c362bbb8e9a12c44953e98eda362582d85ffcab2caa1c8a3c7d0552891f4b10f36b824f49cc3f0abc46cd80f0daa836dc7229795badde0a01009e22068afed18b6f4d64b7c83697fe1722baa98207b9a8165fc52d124986132e61f14154f0adde5eca2cc2664bc510c707a8381b600ff32602cb27d8b5f90c66766859b9a20b3b2d86fbbd4c29cbb8b356957075b74627177ea1ee993bd414014a968ed91b2521bcd05c261aa048c94b71e303ea4a4b5470e13e25eca45a6af274c31eac2f1055f2a87fd65cdaff412f7febb704b3ef2c7e1ccedc3e9edc3547f83b07ad0eac4efc297a8da58ea6c50f45f3aae56529702e397e3d3aa668345f6d3b06eb553b102c3cf5d835e6ec330699176200750ad6c2f3a69bb638fd26b543aedd93eef914ad48292d8f4a55a67cfddbeb4c0c8ea10d5d143a625ec6fe6adf91fc0fab226db1b2a740fd2a9e757e34b31cf175ed061057710cba6e4937e40b1303d3b2689908ea76563de04648ff205cfb1d6b616a7c35ff36af60fa5b9859738163f8cfa3dcd1d11c7cf5b2b9fe5e7e66985a0f5e4f3c51627e92046471d210ba744cdbd53a93ffd4146eec8d9148825388a9ae73ee91a061cc5587e60c5ea67b50ad63bc418fe1d447a484860cb7bf191a807f660399044e037971ff15d3b142abd1e81006b5fe5aa4813b6814695ab54c393e91c3049e4edea9c369329369f15f0ad386e2a483931d0e31357ea81cd8a6c36a27ba5956249ca4d18918beccd0b11d8a9c52a74ce3a860b0c6d8d0158273c4e9e2835936b51b0e30aafca53f5dc5b456c1d9887f4c70c97cb7bfd4f26eae75298acc62274ae9048ba9fc96b895911371138237b3c4374ff071c9cb75a774812a24c3f278f96c36625711e91dffffa81a7cd25ea0dd7c6c1d7a6973e48d57a3ce089d9f99ad5e44a1f4210521aec8d880e667954ed663251966fd8b000618f2ffafac873f6d7b97ff1304ed279a873721c56eb5e1dc0e8d51667a73415f483226572772749f85fc126754054c664c46242d608aa759a543141ad35d424746a5d1be9fa05273cc3398c11d16aeb2ed81d4b059b08b8f02fd1ded5a6c0ad72a02a8cdccdffbf8f2b196278dad14ff4545ab95146e394bd5e8ea34df418a047140a132590889603e16d9c83f3f6aa55e1330caae59063ba57fdd9c11a1ec4c4bb141a40724f5f7601ff8bc4ad3a9c9f563e45a29fb8c74d0c13934f83db1b5a52a8d03055ead3245fd727ee45058330e91e712a5380ab656a50b4b5820fb3b49d4e84e0de71b7452a427f6ff7350eb3f952d45904ffc250762738604062e5a097fca5d21ded3077d1f30201177ebecfe92bed88a1ac561d9fefc904d8a1b09a5e20682a6ccc59be4934f24e2480ede28cd62cdac7b68c9ce27cd333fbb072fb8b5ce5a8735d9ea154605b91e0668d68354c55f581bdcb69f250a7e56487c1ba809fb6c7f7ceee5907c91021dbe24a662c12ee3d6ff4633b4954e88a40d60cfee979ddfcb8116b01f22f4b11d1b5eecd7561d4b16f83d04b4814c28882f6068fb949fd09268f59b6abde55050fd776b1b424dfd2802556aaede8b1afa56720c8d13b05c2189d72f1e3708c174f8fb688462dfdc70e84b44f7da3d616efa1ead9ef1c4337699879bad47f79a6e5576a3814453091d80ed9d03cd4216b4c143212576efd0e836abb499cfa3f9e602c919714cb448adae0339278f7f1258ce23155a26e2fa0e6596aff470bf1086a93b63960eef567a28750f858ead93ecca9b7ca29a842e463a871df69756ecf670b76415a0867794755a179a88128b201f332edc6cd08b6a7beafe75b0821bd7b1022735c2bfd3f90af4e36e847fd6b97553d9560a145b74c27a31140b48f9f308d019081ac936e904d70a3356d793e639f7792d008f065279a4a28c4809b85c05526b59915eb46dc8ac5653347648f7fbc5a667a3d5944087fd9efaff19fe13570d24c2beda45e7644bd5497e195dc2e48814103406db0785d619ae53086fce664af805d7e61f3298d565a7237ab0629a190a48c9bed81dbac82e1b5acc11fc9ff2c9c8248a5bb335e894b72bbddc73189c24cab836dc1f6c8128e9133e2b27d20766100a0b98f6b91839fbe1fb154a6cdd5be0b6c13f78f87823f8c55169779eb1b117c20846b03ed2bac2b18ac74e2b64c3f82759c581e5e53154dbade46765845c033463ff99fe6e078ca0f1aead87493dc482cbd39e4c2ca19edb751dadcab913c478803d69bc1d944b331ba420fcc86ac5390650e495b89b5c6872b42b3f2e16cdf1b6f963e8e0057557dbc2b24e0bcd36387087b3023cab9a2ece6afbd51e83bcd1007b9757326e263055ea1a2ec3125a8af54b02895d1125c0d1ede0cfbf4866171e25fad720a4a50363ad5bfe38e14adda5b5c82e6181d31d30fd23a38ecfe98ee68be01b4a7d9490d18d7425aee70d06e4553f149832c8851a0dfeee41ba9b0e4f7ab46ae2560e5ad6fc35288bccb340d41dc7def1cb3a7f398ac2d1693d543bf98239e5db0e2f34fc14cd8b44d1ba53833620deea78a7d90cb359142d5a61e29d38602ae494ac779e9ec16b7792b757713bd19397cd547f969848d086af521d829a271de1b2a382132daabb9086ee0654160f2e7fb020b10c6ea26b763d20b49a97c5970df8051da2d71a32f0465ecae195d1b70b08bf23017c27ffd2cc897dbfb84cc9ec7a8dd177d8591a6aaf77dc306ca2328f67c53d48283e6233ad61ef542f5b6e130c012b7c9e69a61b121168ba8863ff45311ad9561cb4b5a5c9ed9745a37ab352d93ab88482e01652fcce357bf02c5210399984bd8faf5ebda5f483600186e8ae84698f9b463e021fbf5da2412c156facb8097e14e18ff00fc75acee475bc1285f0ab3c4dc251020e52913f9e1124b14e04c26fa34b9baa11bd0886f3650e0c390d8ecfcac645849699415ffd875afa2b9b7c6b2c7f78491b1faae4778ad6cb47b35c1bfbb3062c61459df8dc034a2f5e4dc12b56c0667f5f3cd7e2b03d3d157687a60eb0063455403516e1fed3890f862a2017366628bf76c85d7b70d75366c83d847f44ce4e51edf158ebde00f99e72d361c2ea2ad04aa8cee7e2f0eff852027ca8339794b4aeb96e52da49e4fe2a98a3b7504118c737c11afec36f1c30a3e9b6a36d5335a9f6b6f242e710d2ea3ab2a468c9b8da4b84470a558ea7ba18f99d517037db2d322449a2d589e0d2e9f93ac94170da72f643c8c647db1e035bde28001c869c8a5258290501b1630b14eac412809767670cb440085ddef76c7d43591967ca3bafaeb177dc065f87a9212fad6481121996a3f5a1e51519d5eb4cea3cb274fc4fa1c5667cbb617fb14f4367c24f09666970172f500a2270ddffd7697aaa096cb6bd47406a7d277c7fb34b887465ede108e803fe7b8410b5ea7577ad8a604d6b7ecdf6cf276c99b045216424c20a7c3da4c9c123e94ef8ee4c0299f6d2622b601167b60ef42bf931275b17a9ef2da9579b743ce087ac79245fab54fea7979c7f2e4e2c8d6758066fa1c6876e3103c908cffe5ca92b58e473d8739d199402ede53cca1de17d2f98f6bbe100f1277ebe7d83117d2c43f600977de9bf30970a5fbc2899e5b089926275760cf3668fd758310a9c74e3fc141316b0b045b2e16b84cae521c0ad8f94926baa5a20bb0c71f467b8c61bfc16f1cd20b5870803e4205cfd601c2c08aa7752a101019bd067844dfb60486edf68ee68ba0b5eae49717fe018e2bf6ef2b67015353dad152f7cb9a524e716d80ac54c274e9e8f4b84e8b0feec2c968c578f20277bea85fcb752218bb1bcea2a6a4dd90367ccb8a2b4fd7a7d3a04b0e1d1816cff222be77a2b525fd"
)

type AK struct {
	HostUUID   string    `json:"host_uuid"`
	MacAddress string    `json:"mac_address"`
	CtlURI     string    `json:"ctl_uri"`
	Expiry     time.Time `json:"expiry"`
	VerifyKey  string    `json:"verify_key"`
}

func AkExist() bool {
	_, err := os.Stat(ConfFile())
	return err == nil
}

func LoadAK() (*AK, error) {
	akBytes, err := os.ReadFile(ConfFile())
	if err != nil {
		return nil, fmt.Errorf("cannot read activation key file: %s\n", err)
	}
	content, err := helper.DecryptMessageArmored(decrypt(k, a, i), nil, string(akBytes[:]))
	if err != nil {
		return nil, fmt.Errorf("invalid activation key: %s\n", err)
	}
	ak := new(AK)
	err = json.Unmarshal([]byte(content), ak)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal activation key: %s\n", err)
	}
	return ak, nil
}

func decrypt(key string, ct string, iv string) string {
	keyBytes, _ := hex.DecodeString(key)
	ciphertext, _ := hex.DecodeString(ct)
	ivBytes, _ := hex.DecodeString(iv)
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		panic(err.Error())
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	plaintext, err := aesgcm.Open(nil, ivBytes, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	s := string(plaintext[:])
	return s
}
