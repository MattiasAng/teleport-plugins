module github.com/gravitational/teleport-plugins

go 1.16

require (
	github.com/alecthomas/kong v0.2.17
	github.com/dgraph-io/badger/v3 v3.2103.2
	github.com/ghodss/yaml v1.0.0
	github.com/go-resty/resty/v2 v2.3.0
	github.com/gogo/protobuf v1.3.2
	github.com/google/btree v1.0.1 // indirect
	github.com/google/go-querystring v1.0.0
	github.com/google/uuid v1.3.0
	github.com/gravitational/kingpin v2.1.11-0.20190130013101-742f2714c145+incompatible
	github.com/gravitational/teleport-plugin-framework v0.0.0-00010101000000-000000000000
	github.com/gravitational/teleport/api v0.0.0-20220330084532-4fce1a4ff991 // tag v9.0.1
	github.com/gravitational/trace v1.1.18
	github.com/hashicorp/go-version v1.3.0
	github.com/hashicorp/terraform-plugin-framework v0.6.1
	github.com/hashicorp/terraform-plugin-go v0.8.0
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.10.1
	github.com/jonboulle/clockwork v0.2.2
	github.com/json-iterator/go v1.1.10
	github.com/julienschmidt/httprouter v1.3.0
	github.com/mailgun/holster/v3 v3.15.2
	github.com/mailgun/mailgun-go/v4 v4.5.3
	github.com/manifoldco/promptui v0.8.0
	github.com/pelletier/go-toml v1.8.0
	github.com/peterbourgon/diskv/v3 v3.0.0
	github.com/sethvargo/go-limiter v0.7.2
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.1
	golang.org/x/net v0.0.0-20220127200216-cd36cc0744dd
	golang.org/x/sync v0.0.0-20201020160332-67f06af15bc9
	golang.org/x/tools v0.0.0-20210106214847-113979e3529a // indirect
	google.golang.org/genproto v0.0.0-20210223151946-22b48be4551b // indirect
	google.golang.org/grpc v1.45.0
	google.golang.org/protobuf v1.27.1
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/mail.v2 v2.3.1
	k8s.io/apimachinery v0.20.4
)

replace (
	github.com/gogo/protobuf => github.com/gravitational/protobuf v1.3.2-0.20201123192827-2b9fcfaffcbf
	github.com/gravitational/teleport-plugin-framework => github.com/gzigzigzeo/teleport-plugin-framework v0.0.0-20220405141808-a80bd9d8a6c0
	github.com/julienschmidt/httprouter => github.com/rw-access/httprouter v1.3.1-0.20210321233808-98e93175c124
)
