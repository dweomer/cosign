//
// Copyright 2021 The Sigstore Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cli

import (
	"fmt"
	"os"

	cranecmd "github.com/google/go-containerregistry/cmd/crane/cmd"
	"github.com/google/go-containerregistry/pkg/logs"
	"github.com/sigstore/cosign/v2/cmd/cosign/cli/options"
	"github.com/sigstore/cosign/v2/cmd/cosign/cli/templates"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	cobracompletefig "github.com/withfig/autocomplete-tools/integrations/cobra"
	"sigs.k8s.io/release-utils/version"
)

var (
	ro = &options.RootOptions{}
)

func normalizeCertificateFlags(_ *pflag.FlagSet, name string) pflag.NormalizedName {
	switch name {
	case "cert":
		name = "certificate"
	case "cert-email":
		name = "certificate-email"
	case "cert-chain":
		name = "certificate-chain"
	case "cert-oidc-issuer":
		name = "certificate-oidc-issuer"
	case "output-cert":
		name = "output-certificate"
	case "cert-identity":
		name = "certificate-identity"
	}
	return pflag.NormalizedName(name)
}

func New() *cobra.Command {
	var (
		out, stdout *os.File
	)

	cmd := &cobra.Command{
		Use:               "cosign",
		Short:             "A tool for Container Signing, Verification and Storage in an OCI registry.",
		DisableAutoGenTag: true,
		SilenceUsage:      true, // Don't show usage on errors
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			if ro.OutputFile != "" {
				var err error
				out, err = os.Create(ro.OutputFile)
				if err != nil {
					return fmt.Errorf("error creating output file %s: %w", ro.OutputFile, err)
				}
				stdout = os.Stdout
				os.Stdout = out // TODO: don't do this.
				cmd.SetOut(out)
			}

			if ro.Verbose {
				logs.Debug.SetOutput(os.Stderr)
			}

			return nil
		},
		PersistentPostRun: func(_ *cobra.Command, _ []string) {
			if out != nil {
				_ = out.Close()
			}
			os.Stdout = stdout
		},
	}
	ro.AddFlags(cmd)

	templates.SetCustomUsageFunc(cmd)

	// Add sub-commands.
	cmd.AddCommand(Attach())
	cmd.AddCommand(Attest())
	cmd.AddCommand(AttestBlob())
	cmd.AddCommand(Bundle())
	cmd.AddCommand(Clean())
	cmd.AddCommand(Debug())
	cmd.AddCommand(Tree())
	cmd.AddCommand(Completion())
	cmd.AddCommand(Copy())
	cmd.AddCommand(Dockerfile())
	cmd.AddCommand(Download())
	cmd.AddCommand(Generate())
	cmd.AddCommand(GenerateKeyPair())
	cmd.AddCommand(ImportKeyPair())
	cmd.AddCommand(Initialize())
	cmd.AddCommand(Load())
	cmd.AddCommand(Manifest())
	cmd.AddCommand(PIVTool())
	cmd.AddCommand(PKCS11Tool())
	cmd.AddCommand(PublicKey())
	cmd.AddCommand(Save())
	cmd.AddCommand(Sign())
	cmd.AddCommand(SignBlob())
	cmd.AddCommand(Upload())
	cmd.AddCommand(Verify())
	cmd.AddCommand(VerifyAttestation())
	cmd.AddCommand(VerifyBlob())
	cmd.AddCommand(VerifyBlobAttestation())
	cmd.AddCommand(Triangulate())
	cmd.AddCommand(TrustedRoot())
	cmd.AddCommand(SigningConfig())
	cmd.AddCommand(Env())
	cmd.AddCommand(version.WithFont("starwars"))

	cmd.AddCommand(cranecmd.NewCmdAuthLogin("cosign"))

	cmd.SetGlobalNormalizationFunc(normalizeCertificateFlags)
	cmd.AddCommand(cobracompletefig.CreateCompletionSpecCommand())

	return cmd
}
