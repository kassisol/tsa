package tenant

import (
	"os"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/kassisol/tsa/cli/session"
	"github.com/kassisol/tsa/client"
	tntpkg "github.com/kassisol/tsa/pkg/tenant"
	"github.com/spf13/cobra"
)

var (
	tenantMemberAdd    bool
	tenantMemberRemove bool
)

func newMemberCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "member [tenant_id] [group]",
		Short: "Manage tenant membership to group",
		Long:  memberDescription,
		Run:   runMember,
	}

	flags := cmd.Flags()
	flags.BoolVarP(&tenantMemberAdd, "add", "a", false, "Add user to group")
	flags.BoolVarP(&tenantMemberRemove, "remove", "r", false, "Remove user from group")

	return cmd
}

func runMember(cmd *cobra.Command, args []string) {
	if len(args) < 2 || len(args) > 2 {
		cmd.Usage()
		os.Exit(-1)
	}

	sess, err := session.New()
	if err != nil {
		log.Fatal(err)
	}
	defer sess.End()

	srv, err := sess.Get()
	if err != nil {
		log.Fatal(err)
	}

	clt, err := client.New(srv.Server.TSAURL)
	if err != nil {
		log.Fatal(err)
	}

	tid, err := strconv.Atoi(args[0])
	if err != nil {
		log.Fatal(err)
	}

	itnt, err := tntpkg.New()
	if err != nil {
		log.Fatal(err)
	}

	if err := itnt.SetTenant(tid); err != nil {
		log.Fatal(err)
	}

	tenant2 := itnt.GetTenant()

	if tenantMemberAdd {
		if err := clt.AddGroupToTenant(srv.Token, tid, tenant2.Name, args[1]); err != nil {
			log.Fatal(err)
		}
	}
	if tenantMemberRemove {
		var groupe int
		for _, grp := range tenant2.AuthGroups {
			if args[1] == grp.Name {
				groupe = int(grp.ID)

				break
			}
		}

		if err := clt.RemoveGroupFromTenant(srv.Token, tid, groupe); err != nil {
			log.Fatal(err)
		}
	}
}

var memberDescription = `
Manage tenant membership to group

`
