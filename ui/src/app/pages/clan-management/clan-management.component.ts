import {Component, OnDestroy, OnInit, ViewChild} from '@angular/core';
import {AuthService} from "../../services/auth.service";
import {FormBuilder} from "@angular/forms";
import {ClansService} from "../../services/clans.service";
import {Subject} from "rxjs";
import {BasicStats, ClanInfo, ClanInvite, Player, PlayerClaims} from "../../models";
import {takeUntil} from "rxjs/operators";
import {NgbModal} from "@ng-bootstrap/ng-bootstrap";

@Component({
  selector: 'app-clan-management',
  templateUrl: './clan-management.component.html',
  styleUrls: ['./clan-management.component.sass']
})
export class ClanManagementComponent implements OnInit, OnDestroy {

  @ViewChild('registerCongrats') private registerCongrats;
  @ViewChild('leaveSure') private leaveSure;
  @ViewChild('kickSure') private kickSure;

  componentDestroyed$ = new Subject();

  registerClanForm = this.formBuilder.group({
    clanname: ''
  });
  inviteMemberForm = this.formBuilder.group({
    playername: ''
  });
  hideRegisterError = true;
  hideInviteError = true;

  playerClaims: PlayerClaims = null;
  clan: ClanInfo = null;
  disbandModal = null;
  leaveModal = null;
  kickModal = null;
  clanInvites: ClanInvite[];
  myInvites: ClanInvite[];
  members: BasicStats[];
  kickSurePlayer: Player;

  constructor(private formBuilder: FormBuilder,
              private clansService: ClansService,
              private authService: AuthService,
              private modalService: NgbModal) { }

  ngOnInit() {
    this.authService.playerClaims.pipe(takeUntil(this.componentDestroyed$))
      .subscribe( value => {
        this.playerClaims = value;
        if (this.playerClaims && this.playerClaims.clanID != null && this.playerClaims.clanID !== 0) {
          this.clansService.getClan(this.playerClaims.clanID).subscribe( clan => {
            this.clan = clan;
            this.updateMembers();
            if (this.clan.leaderID === this.playerClaims.playerID) {
              this.updateClanInvites();
            }
          });
        } else {
          this.updateMyClanInvites();
        }
      });
  }

  ngOnDestroy(): void {
    this.componentDestroyed$.next();
  }

  openDisband(content): void {
    this.disbandModal = this.modalService.open(content, {ariaLabelledBy: 'modal-basic-title'});
  }

  openCongrats(content): void {
    this.modalService.open(content, {ariaLabelledBy: 'modal-basic-title'});
  }

  disbandClan(): void {
    this.disbandModal.dismiss();
    this.clansService.disband(this.clan.id).subscribe(resp => {
      this.clan = null;
      this.authService.revalidate();
    });
  }

  onRegister(): void {
    const clanname = this.registerClanForm.value.clanname;
    this.clansService.register(clanname).subscribe( resp => {
      this.authService.revalidate();
      this.hideRegisterError = true;
      this.registerClanForm.reset();
      this.clan = resp;
      this.openCongrats(this.registerCongrats);
    }, err => {
      this.hideRegisterError = false;
    });
  }

  updateClanInvites(): void {
    this.clansService.getClanInvites(this.clan.id).subscribe( resp => {
        this.clanInvites = resp;
    });
  }

  updateMyClanInvites(): void {
    this.clansService.getMyInvites().subscribe( resp => {
      this.myInvites = resp;
    });
  }

  cancelInvite(playerID: number): void {
    this.clansService.cancelClanInvite(this.clan.id, playerID).subscribe( resp => {
      this.updateClanInvites();
    });
  }

  onInviteMember(): void {
    const playername = this.inviteMemberForm.value.playername;
    this.clansService.invite(this.clan.id, playername).subscribe( resp => {
      this.hideInviteError = true;
      this.updateClanInvites();
      this.inviteMemberForm.reset();
    }, err => {
      this.hideInviteError = false;
    });
  }

  updateMembers(): void {
    this.clansService.getMembers(this.clan.id).subscribe( resp => {
      this.members = resp;
    });
  }

  declineInvite(clanID: number): void {
    this.clansService.declineClanInvite(clanID).subscribe( resp => {
      this.updateMyClanInvites();
    });
  }

  acceptInvite(clanID: number): void {
    this.clansService.acceptClanInvite(clanID).subscribe(resp => {
      this.authService.revalidate();
    }, err => {
      this.updateMyClanInvites();
    });
  }

  openKick(player): void {
    this.kickSurePlayer = player;
    this.kickModal = this.modalService.open(this.kickSure, {ariaLabelledBy: 'modal-basic-title'});
  }

  kickPlayer(playerID: number): void {
    this.kickModal.dismiss();
    this.clansService.kick(this.clan.id, playerID).subscribe(resp => {
      this.updateMembers();
    });
  }

  openLeave(): void {
    this.leaveModal = this.modalService.open(this.leaveSure, {ariaLabelledBy: 'modal-basic-title'});
  }

  leaveClan(): void {
    this.leaveModal.dismiss();
    this.clansService.leave().subscribe(resp => {
      this.clan = null;
      this.authService.revalidate();
    });
  }

}
