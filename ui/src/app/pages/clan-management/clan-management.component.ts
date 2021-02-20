import {Component, OnInit, ViewChild} from '@angular/core';
import {AuthService} from "../../services/auth.service";
import {FormBuilder} from "@angular/forms";
import {ClansService} from "../../services/clans.service";
import {Subject} from "rxjs";
import {ClanInfo, PlayerClaims} from "../../models";
import {takeUntil} from "rxjs/operators";
import {NgbModal} from "@ng-bootstrap/ng-bootstrap";

@Component({
  selector: 'app-clan-management',
  templateUrl: './clan-management.component.html',
  styleUrls: ['./clan-management.component.sass']
})
export class ClanManagementComponent implements OnInit {

  @ViewChild('registerCongrats') private registerCongrats;

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
  modal = null;

  constructor(private formBuilder: FormBuilder,
              private clansService: ClansService,
              private authService: AuthService,
              private modalService: NgbModal) { }

  ngOnInit() {
    this.authService.playerClaims.pipe(takeUntil(this.componentDestroyed$))
      .subscribe( value => {
        this.playerClaims = value;
        console.log('value');
        console.log(value);
        if (this.playerClaims && this.playerClaims.clanID != null && this.playerClaims.clanID !== 0) {
          this.clansService.getClan(this.playerClaims.clanID).subscribe( clan => {
            this.clan = clan;
          });
        }
      });
  }

  ngOnDestroy(): void {
    this.componentDestroyed$.next();
  }

  openDisband(content): void {
    this.modal = this.modalService.open(content, {ariaLabelledBy: 'modal-basic-title'});
  }

  openCongrats(content): void {
    this.modalService.open(content, {ariaLabelledBy: 'modal-basic-title'});
  }

  disbandClan(): void {
    this.clansService.disband(this.clan.id).subscribe(resp => {
      this.modal.dismiss();
      this.clan = null;
      this.authService.revalidate();
    }, err => {
      this.modal.dismiss();
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

  onInviteMember(): void {
    const playername = this.inviteMemberForm.value.playername;
    console.log('invite player: ' + playername);
  }

}
